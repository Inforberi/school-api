package middlewares

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h := w.Header()

		// Почти всегда безопасно и полезно.
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin") // хороший баланс (не ломает аналитику так сильно)
		h.Set("X-Frame-Options", "DENY")                            // legacy; реальная защита ниже через CSP frame-ancestors
		h.Set("X-DNS-Prefetch-Control", "off")
		h.Set("X-Permitted-Cross-Domain-Policies", "none")
		h.Set("X-XSS-Protection", "0") // устарело — выключаем

		// Не раскрываем стек.
		h.Del("X-Powered-By")

		// CSP: МЯГКИЙ старт, чтобы не ломать CDN/аналитику/шрифты.
		// Это не "максимальная" защита, но обычно не ломает типичный фронт.
		// TODO: ужесточать по мере понимания ресурсов (лучший путь — собрать список доменов и ограничить).
		h.Set("Content-Security-Policy",
			"default-src 'self' https: data: blob:; "+
				"base-uri 'self'; "+
				"object-src 'none'; "+
				"frame-ancestors 'none'; "+
				"form-action 'self'; "+
				"upgrade-insecure-requests",
		)

		// COOP: оставляем "same-origin-allow-popups", чтобы не ломать window.open/SSO-потоки.
		h.Set("Cross-Origin-Opener-Policy", "same-origin-allow-popups")

		// COEP: НЕ включаем, чтобы не ломать загрузку ресурсов без CORP/CORS.
		// h.Set("Cross-Origin-Embedder-Policy", "require-corp")

		// CORP: тоже может ломать выдачу ресурсов для других origin, поэтому не ставим глобально.
		// h.Set("Cross-Origin-Resource-Policy", "same-origin")

		// Permissions-Policy: мягкий дефолт (запрещаем чувствительное).
		// Обычно не ломает сайт, если вы не используете эти API.
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// HSTS: включаем только при прямом TLS на Go, чтобы не получить неверное поведение за прокси.
		if r.TLS != nil {
			h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		next.ServeHTTP(w, r)
	})
}
