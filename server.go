package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/clbanning/mxj"
)

var	desiredQuality = "FLASH_450K_400X224" 
// from 450K (50MB), _1200K_640X360 (150MB), _1800K_960X540 (200MB) or _2500K_1280X720 (300MB)
				
func main() {
	date := time.Now().AddDate(0, 0, -1)
	dates := "year_" + date.Format("2006/month_01/day_02")

	myTeamsMap := InitMyTeamsMap()
	games := make(map[int][]string)
	games = SearchMyMLBGames(dates, games, myTeamsMap)

	downloadedGames := downloadMyMLBGames(games, myTeamsMap)
	log.Printf("%v", downloadedGames)

	pushBulletAPI := os.Getenv("pushBulletAPI")
	log.Printf(pushBulletAPI)

	// TODO3: prepare upload_urls per game

	// TODO4: upload games to pushbullet

	// TODO5: send file via pushbullet
}


func downloadMyMLBGames(games map[int][]string, myTeamsMap map[int]string) []string {
	var downloadedGames []string
	downloadedGames = make ([]string, len(games))

	for _, v := range games {
		awayTeamID, _ := strconv.Atoi(v[0])
		homeTeamID, _ := strconv.Atoi(v[1])
		fileName := "/Users/family/Downloads/" + myTeamsMap[awayTeamID] + "@" + myTeamsMap[homeTeamID] + "_" + desiredQuality + ".mp4"
		fileName = strings.Replace(fileName, " ", "+", -1)
		out, _ := os.Create(fileName)
		defer out.Close()
		resp, _ := http.Get(v[2])
		defer resp.Body.Close()
		_, err := io.Copy(out, resp.Body)

		if err == nil {
			downloadedGames = append(downloadedGames, fileName)
		}
	}

	return downloadedGames
}
// SearchMLBGames is now commented
func SearchMyMLBGames(date string, games map[int][]string, myTeamsMap map[int]string) map[int][]string {
	domain := "http://gd2.mlb.com/components/game/mlb/"
	suffix := "/grid_ce.xml"
	url := domain + date + suffix
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	xml, err := ioutil.ReadAll(resp.Body)
	m, err := mxj.NewMapXml(xml)

	gameInfos, err := m.ValuesForKey("game")
	if err != nil {
		log.Fatal("err:", err.Error())
		log.Printf("MLB site '%s' response empty", domain)
		games[0] = []string{"Error connecting to " + domain}
		return games
	}

	// now just manipulate Map entries returned as []interface{} array.
	for k, v := range gameInfos {
		gameID := ""
		aGameVal, _ := v.(map[string]interface{})
		if aGameVal["-media_state"].(string) == "media_dead" {
			continue
		}

		// rescan looking for keys with data: Values or Value
		gm := aGameVal["game_media"].(map[string]interface{})
		hb := gm["homebase"].(map[string]interface{})
		media := hb["media"].([]interface{})
		for _, val := range media {
			aMediaVal, _ := val.(map[string]interface{})
			if aMediaVal["-type"].(string) != "condensed_game" {
				continue
			} else {
				gameID = aMediaVal["-id"].(string)
				continue
			}
		}

		if gameID != "" {
			fetchGame := false

			awayTeamIDs := aGameVal["-away_team_id"].(string)
			awayTeamID, _ := strconv.Atoi(awayTeamIDs)
			homeTeamIDs := aGameVal["-home_team_id"].(string)
			homeTeamID, _ := strconv.Atoi(homeTeamIDs)
			if myTeamsMap[awayTeamID] != "" {
				fetchGame = true
			} else if myTeamsMap[homeTeamID] != "" {
				fetchGame = true
			}

			if fetchGame {
				// grab mp4 URLs per game
				detailURL := "http://m.mlb.com/gen/multimedia/detail" + generateDetailURL(gameID)
				gameURL := fetchGameURL(detailURL, desiredQuality)

				games[k] = []string{awayTeamIDs, homeTeamIDs, gameURL}
			}
		}
	}

	return games
}


// fetchGameURL is now commented
func fetchGameURL(detailURL string, desiredQuality string) string {
	gameURL := "MickeyMouse.mp4"

	resp, err := http.Get(detailURL)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	xml, err := ioutil.ReadAll(resp.Body)
	m, err := mxj.NewMapXml(xml)

	URLs, err := m.ValuesForKey("url")
	if err != nil {
		log.Fatal("err:", err.Error())
		return ""
	}

	// now just manipulate Map entries returned as []interface{} array.
	for _, v := range URLs {
		aGameVal, _ := v.(map[string]interface{})
		if aGameVal["-playback_scenario"].(string) == desiredQuality {
			return aGameVal["#text"].(string)
		}
	}

	return gameURL
}


// generateDetailURL is now commented
func generateDetailURL(gameID string) string {
	// given gameID 605442983 return "/9/8/3/605442983.xml"
	return 	"/" + gameID[len(gameID)-3:len(gameID)-2] + 
					"/" + gameID[len(gameID)-2:len(gameID)-1] + 
					"/" + gameID[len(gameID)-1:] + 
					"/" + gameID + ".xml"
}


// InitMyTeamsMap is now commented
func InitMyTeamsMap() map[int]string {
	myTeamsMap := make(map[int]string)
  myTeamsMap[111] = "Boston Red Sox"
	myTeamsMap[139] = "Tampa Bay Rays"
  myTeamsMap[141] = "Toronto Blue Jays"
  myTeamsMap[120] = "Washington Nationals"
  myTeamsMap[137] = "San Francisco Giants"

/*	homePageMap[110] = Team{"Baltimore Orioles", "http://m.orioles.mlb.com/roster"}
	homePageMap[145] = Team{"Chicago Whitesox", "http://m.whitesox.mlb.com/roster"}
	homePageMap[117] = Team{"Houston Astros", "http://m.astros.mlb.com/roster"}
	homePageMap[144] = Team{"Atlanta Braves", "http://m.braves.mlb.com/roster"}
	homePageMap[112] = Team{"Chicago Cubs", "http://m.cubs.mlb.com/roster"}
	homePageMap[109] = Team{"Arizona Diamond Backs", "http://m.dbacks.mlb.com/roster"}
	homePageMap[111] = Team{"Boston Red Sox", "http://m.redsox.mlb.com/roster"}
	homePageMap[114] = Team{"Cleveland Indians", "http://m.indians.mlb.com/roster"}
	homePageMap[108] = Team{"Los Angeles Angels", "http://m.angels.mlb.com/roster"}
	homePageMap[146] = Team{"Miami Marlins", "http://m.marlins.mlb.com/roster"}
	homePageMap[113] = Team{"Cincinnati Reds", "http://m.reds.mlb.com/roster"}
	homePageMap[115] = Team{"Colorado Rockies", "http://www.rockies.com/roster"}
	homePageMap[147] = Team{"New York Yankees", "http://m.yankees.mlb.com/roster"}
	homePageMap[116] = Team{"Detroit Tigers", "http://www.tigers.com/roster"}
	homePageMap[133] = Team{"Oakland Athletics", "http://m.athletics.mlb.com/roster"}
	homePageMap[121] = Team{"New York Mets", "http://m.mets.mlb.com/roster"}
	homePageMap[158] = Team{"Milwaukee Brewers", "http://m.brewers.mlb.com/roster"}
	homePageMap[119] = Team{"LA Dodgers", "http://m.dodgers.mlb.com/roster"}
	homePageMap[139] = Team{"Tampa Bay Rays", "http://m.rays.mlb.com/roster"}
	homePageMap[118] = Team{"Kansas City Royals", "http://m.royals.mlb.com/roster"}
	homePageMap[136] = Team{"Seattle Mariners", "http://m.mariners.mlb.com/roster"}
	homePageMap[143] = Team{"Philadelphia Phillies", "http://m.phillies.mlb.com/roster"}
	homePageMap[138] = Team{"St Louis Cardinals", "http://m.cardinals.mlb.com/roster"}
	homePageMap[135] = Team{"San Diego Padres", "http://m.padres.mlb.com/roster"}
	homePageMap[141] = Team{"Toronto Blue Jays", "http://m.bluejays.mlb.com/roster"}
	homePageMap[142] = Team{"Minnesota Twins", "http://m.twins.mlb.com/roster"}
	homePageMap[140] = Team{"Texas Rangers", "http://m.rangers.mlb.com/roster"}
	homePageMap[120] = Team{"Washington Nationals", "http://m.nationals.mlb.com/roster"}
	homePageMap[134] = Team{"Pittsburgh Pirates", "http://m.pirates.mlb.com/roster"}
	homePageMap[137] = Team{"San Francisco Giants", "http://m.giants.mlb.com/roster"}*/

	return myTeamsMap
}