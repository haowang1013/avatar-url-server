package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/op/go-logging"
	"net/http"
	"os"
)

const (
	port = 8080
)

var (
	steamAPIKey string
	log         = logging.MustGetLogger("")
)

/**
{
  "response": {
    "players": [
      {
        "avatar": "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb.jpg",
        "avatarfull": "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb_full.jpg",
        "avatarmedium": "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb_medium.jpg",
        "communityvisibilitystate": 3,
        "lastlogoff": 1482849808,
        "personaname": "wang hao",
        "personastate": 1,
        "personastateflags": 0,
        "primaryclanid": "103582791430123379",
        "profilestate": 1,
        "profileurl": "http://steamcommunity.com/profiles/76561197968196788/",
        "realname": "wanghao",
        "steamid": "76561197968196788",
        "timecreated": 1092512230
      }
    ]
  }
}
*/

type SteamPlayersResponse struct {
	Response SteamPlayers `json:"response"`
}

type SteamPlayers struct {
	Players []SteamPlayerSummary `json:"players"`
}

type SteamPlayerSummary struct {
	Avatar       string `json:"avatar"`
	AvatarFull   string `json:"avatarfull"`
	AvatarMedium string `json:"avatarmedium"`
	ProfileUrl   string `json:"profileurl"`
	PersonaName  string `json:"personaname"`
	RealName     string `json:"realname"`
	SteamID      string `json:"steamid"`
}

func init() {
	steamAPIKey = os.Getenv("STEAM_API_KEY")
	if len(steamAPIKey) == 0 {
		panic("Failed to get steam api key from env variable 'STEAM_API_KEY'")
	}

	format := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
	)

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	formtter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(formtter)
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/avatar/steam/:steamid", func(c *gin.Context) {
		steamID := c.Param("steamid")
		url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", steamAPIKey, steamID)
		resp, err := grequests.Get(url, nil)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		parsed := new(SteamPlayersResponse)
		err = resp.JSON(parsed)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		if len(parsed.Response.Players) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "player doesn't exist",
			})
			return
		}

		if len(parsed.Response.Players) > 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "more than 1 player is found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"avatar_url": parsed.Response.Players[0].AvatarFull,
		})
	})

	log.Debugf("listen on port %d", port)
	router.Run(fmt.Sprintf(":%d", port))
}
