package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rassakhatsky/tdameritrade/coolStuff"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const (
	authURL = "https://auth.tdameritrade.com/auth?response_type=code&redirect_uri=%s&client_id=%s"

	defaultRedirectURL = "http://localhost:8888/callback"
	defaultClientID    = "GO-AMERITRADE@AMER.OAUTHAP"
)

func init() {
	oauth2.RegisterBrokenAuthHeaderProvider("https://auth.tdameritrade.com")
	oauth2.RegisterBrokenAuthHeaderProvider("https://api.tdameritrade.com")
}

var TDConnectEndpoint = oauth2.Endpoint{
	AuthURL:  "https://auth.tdameritrade.com/auth",
	TokenURL: "https://api.tdameritrade.com/v1/oauth2/token",
}

var oauthCode string

func parseInput(r *bufio.Reader, defValue string) string {
	value, err := r.ReadString('\n')
	if err != nil {
		log.Panic().Err(err).Msg("failed to parse value")
	}
	if len(value) <= 1 {
		return defValue
	}

	return value
}

func RequestToken(timeout int) {
	reader := bufio.NewReader(os.Stdin)

	var (
		clientId, redirectURL string
	)

	fmt.Println("Enter Callback URL")
	fmt.Println(fmt.Sprintf("Default one is: %s", defaultRedirectURL))

	redirectURL = parseInput(reader, defaultRedirectURL)

	fmt.Println("Enter OAuth User ID (Consumer Key)")
	fmt.Println(fmt.Sprintf("Default one is: %s", defaultClientID))

	clientId = parseInput(reader, defaultClientID)

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer ctxCancel()

	conf := &oauth2.Config{
		ClientID:    clientId,
		Scopes:      []string{},
		Endpoint:    TDConnectEndpoint,
		RedirectURL: redirectURL,
	}

	// Good time to start a server
	cancel := make(chan struct{})
	defer close(cancel)

	go CreateAuthServer(ctx, cancel, ":8888")

	// Generate Auth URL
	authURL := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Use authorization URL to get an access.")
	fmt.Println("tdameritrade enforce https for redirect url, so there is only self-generated certs.")
	fmt.Println(authURL)
	fmt.Println()

	// wait till token is requested or timeout
	select {
	case <-ctx.Done():
		fmt.Println("timeout has been reached")
		return
	case <-cancel:
	}

	fmt.Println("Working...")
	fmt.Println()

	fmt.Println(oauthCode)

	// oauth2 library devs leaving in a different world where everyone using broken library
	token, err := conf.Exchange(ctx, oauthCode,
		oauth2.SetAuthURLParam("redirect_uri", redirectURL),
		oauth2.SetAuthURLParam("client_id", clientId),
		oauth2.SetAuthURLParam("refresh_token", ""),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	)

	if err != nil {
		fmt.Println("failed to receive refresh token")
		fmt.Println(err.Error())
		return
	}

	fmt.Println(coolStuff.Exterminate)

	fmt.Println()
	fmt.Println()

	fmt.Println(fmt.Sprintf("Access Token: %s", token.AccessToken))
	fmt.Println(fmt.Sprintf("Refresh Token: %s", token.RefreshToken))
	fmt.Println(fmt.Sprintf("Expiry: %s", token.Expiry.String()))
	fmt.Println(fmt.Sprintf("Type: %s", token.TokenType))
}
