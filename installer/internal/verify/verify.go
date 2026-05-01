package verify

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Result struct {
	OK        bool
	Status    int
	LatencyMs int64
	Body      string
	Hint      string
	Err       error
}

func (r Result) Print(w io.Writer) {
	if r.OK {
		fmt.Fprintf(w, "✅ OK (%d, %d ms) — модель ответила\n", r.Status, r.LatencyMs)
		if r.Body != "" {
			fmt.Fprintf(w, "   ответ: %s\n", trim(r.Body, 200))
		}
		return
	}
	fmt.Fprintf(w, "❌ Связь с GPT не работает\n")
	if r.Status > 0 {
		fmt.Fprintf(w, "   HTTP %d (%d ms)\n", r.Status, r.LatencyMs)
	}
	if r.Err != nil {
		fmt.Fprintf(w, "   ошибка: %v\n", r.Err)
	}
	if r.Body != "" {
		fmt.Fprintf(w, "   тело: %s\n", trim(r.Body, 400))
	}
	if r.Hint != "" {
		fmt.Fprintf(w, "   👉 %s\n", r.Hint)
	}
}

func trim(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func CheckChat(token, baseURL, model string) Result {
	if token == "" {
		return Result{Hint: "Не задан токен"}
	}
	endpoint := strings.TrimRight(baseURL, "/") + "/chat/completions"

	body, _ := json.Marshal(map[string]any{
		"model":      model,
		"max_tokens": 5,
		"messages":   []map[string]string{{"role": "user", "content": "ping"}},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return Result{Err: err, Hint: "некорректный URL: " + endpoint}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "designapi-installer")

	t0 := time.Now()
	resp, err := http.DefaultClient.Do(req)
	lat := time.Since(t0).Milliseconds()

	if err != nil {
		return classifyNetwork(endpoint, lat, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	r := Result{Status: resp.StatusCode, LatencyMs: lat, Body: string(raw)}
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		r.OK = true
	case resp.StatusCode == 401:
		r.Hint = "Токен невалиден или истёк. Сгенерируйте новый в кабинете designapi.ink."
	case resp.StatusCode == 403:
		r.Hint = "Доступ запрещён: проверьте, что токен активен и у него есть доступ к модели " + model + "."
	case resp.StatusCode == 404:
		r.Hint = "Эндпоинт не найден. Убедитесь, что base URL = https://api.designapi.ink/v1"
	case resp.StatusCode == 429:
		r.Hint = "Лимит запросов. Подождите минуту и повторите."
	case resp.StatusCode >= 500:
		r.Hint = "Сервис временно недоступен. Попробуйте позже или напишите в поддержку."
	default:
		r.Hint = fmt.Sprintf("Неожиданный статус %d. Передайте тело ошибки в поддержку.", resp.StatusCode)
	}
	return r
}

func classifyNetwork(endpoint string, lat int64, err error) Result {
	r := Result{Err: err, LatencyMs: lat}
	u, _ := url.Parse(endpoint)
	host := ""
	if u != nil {
		host = u.Hostname()
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		r.Hint = fmt.Sprintf("DNS не резолвит %s. Проверьте интернет / VPN / DNS-сервер. (%s)", host, dnsErr.Err)
		return r
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		r.Hint = fmt.Sprintf("Сетевая ошибка при подключении к %s: %s. Проверьте firewall / прокси.", host, opErr.Op)
		return r
	}
	var tlsErr *tls.CertificateVerificationError
	if errors.As(err, &tlsErr) {
		r.Hint = "TLS не прошёл валидацию сертификата. Возможно, корпоративный MITM-прокси перехватывает трафик."
		return r
	}
	if errors.Is(err, context.DeadlineExceeded) {
		r.Hint = "Таймаут запроса. Сервер не ответил за 15 сек. Проверьте, не блокирует ли вас firewall/прокси."
		return r
	}
	if strings.Contains(err.Error(), "x509") || strings.Contains(err.Error(), "certificate") {
		r.Hint = "Проблема с TLS-сертификатами. Проверьте системное время и корневые сертификаты."
		return r
	}
	r.Hint = "Не удалось установить соединение. Проверьте интернет, прокси и доступность " + host + "."
	return r
}
