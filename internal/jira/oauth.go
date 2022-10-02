package jira

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/andygrunwald/go-jira"
	"github.com/dghubble/oauth1"
	"golang.org/x/net/context"
)

func getJIRAHTTPClient(ctx context.Context, config *oauth1.Config, URL *url.URL) (*http.Client, error) {
	cacheFile, err := jiraTokenCacheFile(URL)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to get path to cached credential file.")
		return nil, err

	}
	tok, err := jiraTokenFromFile(cacheFile)
	if err != nil {
		log.Error().Err(err).Msgf("jiratoken from file ")
		tok = getJIRATokenFromWeb(config)
		saveJIRAToken(cacheFile, tok)
	}
	return config.Client(ctx, tok), nil
}

func getJIRATokenFromWeb(config *oauth1.Config) *oauth1.Token {
	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		log.Fatal().Msgf("Unable to get request token. %v", err)
	}
	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatal().Msgf("Unable to get authorization url. %v", err)
	}
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authorizationURL.String())

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal().Msgf("Unable to read authorization code. %v", err)
	}

	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, code)
	if err != nil {
		log.Fatal().Msgf("Unable to get access token. %v", err)
	}
	return oauth1.NewToken(accessToken, accessSecret)
}

func jiraTokenCacheFile(URL *url.URL) (string, error) {

	return filepath.Join("/conf",
		url.QueryEscape((*URL).Host+".json")), nil
}

func jiraTokenFromFile(file string) (*oauth1.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth1.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveJIRAToken(file string, token *oauth1.Token) {
	log.Info().Msgf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal().Msgf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getJIRAClient(URL *url.URL, ConsumerKey string, jiraPrivateKey string) (*jira.Client, error) {
	ctx := context.Background()
	keyDERBlock, _ := pem.Decode([]byte(jiraPrivateKey))
	if keyDERBlock == nil {
		return nil, errors.New("unable to decode key PEM block")
	}
	if !(keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY")) {
		return nil, fmt.Errorf("unexpected key DER block type: %s", keyDERBlock.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PKCS1 private key. %v", err)
	}
	config := oauth1.Config{
		ConsumerKey: ConsumerKey,
		CallbackURL: "oob", /* for command line usage */
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: (*URL).String() + "plugins/servlet/oauth/request-token",
			AuthorizeURL:    (*URL).String() + "plugins/servlet/oauth/authorize",
			AccessTokenURL:  (*URL).String() + "plugins/servlet/oauth/access-token",
		},
		Signer: &oauth1.RSASigner{
			PrivateKey: privateKey,
		},
	}
	c, err := getJIRAHTTPClient(ctx, &config, URL)
	if err != nil {
		return nil, err
	}
	jiraClient, err := jira.NewClient(c, (*URL).String())
	if err != nil {
		return nil, fmt.Errorf("unable to create new JIRA client. %v", err)
	}
	return jiraClient, nil
}
