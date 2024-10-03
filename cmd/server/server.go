package main

import "effictiveMobile/internal/application"

// @Version 1.0.0
// @Title Song library service
// @Description Сервис представляет из себя библиотеку песен
// @ContactName PaxySong
// @ContactEmail vhser@yandex.ru
// @ContactURL http://someurl.oxox
// @TermsOfServiceUrl http://someurl.oxox
// @LicenseName MIT
// @LicenseURL https://en.wikipedia.org/wiki/MIT_License
// @Server http://localhost:8001 Server-1
// @Security Authorization read write
// @SecurityScheme Authorization apiKey header Authorization
func main() {
	application.Run()
}
