package githupapi

import (
	"fmt"
	"github.com/gittool/base"
	"github.com/gittool/models"
	"github.com/go-resty/resty/v2"
	"regexp"
	"strconv"
	"sync"
)

const DEBUG = true

func GetUserInfo(username string, ch chan *models.Userinfo) {
	token := base.Token()
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	client := resty.New()

	var result *models.Userinfo
	// 最简单的方法
	resp, err := client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetQueryParams(map[string]string{
			"per_page":  "30",
			"page":      "1",
			"sort":      "forks",
			"direction": "des",
		}).
		SetResult(&result).
		Get(url)

	if err != nil {
		ch <- result
	} else {
		_ = resp
		ch <- result
	}
}

func GetUserBaseInfos(username string) *models.Userinfo {

	ch := make(chan *models.Userinfo)
	go GetUserInfo(username, ch)
	result := <-ch
	if DEBUG {
		fmt.Printf("%+v", result)
	}
	return result
}

func GetUserRepos(username string) []*models.Repo {
	userinfo := GetUserBaseInfos(username)
	//fmt.Printf("%+v", userinfo)
	totalPage := userinfo.PublicRepos/100 + 1
	repo_list := make([]*models.Repo, 0)
	for i := 1; i <= totalPage; i++ {
		repo_list = append(repo_list, GetUserRepo(username, totalPage)...)
	}
	//fmt.Println(repo_list)
	return repo_list
}
func GetUserRepo(username string, page int) []*models.Repo {
	token := base.Token()
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=500", username)

	client := resty.New()

	var result []*models.Repo
	// 最简单的方法
	resp, err := client.R().
		SetAuthToken(token).
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetQueryParams(map[string]string{
			"per_page":  "30",
			"page":      strconv.Itoa(page),
			"sort":      "forks",
			"direction": "des",
		}).
		SetResult(&result).
		Get(url)

	base.Must(err)
	//fmt.Println(resp)
	_ = resp
	repo_list := make([]*models.Repo, 0)

	for _, repo := range result {
		if repo.Fork == false {
			//fmt.Println(repo.Name)
			repo_list = append(repo_list, repo)
			fmt.Println(repo.Name)
			fmt.Println(repo.Description)
			fmt.Println()
		}
	}
	//fmt.Println(len(repo_list))
	return repo_list
}

func getFollower(username string, count int) []string {
	var followers []*models.Follower
	var followers_list []string

	var total_page int
	total_page = count/100 + 1
	url := fmt.Sprintf("https://api.github.com/users/%s/followers?per_page=500", username)

	client := resty.New()

	for i := 1; i <= total_page; i++ {
		// 最简单的方法
		resp, err := client.R().
			SetHeader("Accept", "application/vnd.github.v3+json").
			SetResult(&followers).
			SetQueryParams(map[string]string{
				"per_page":  "100",
				"page":      strconv.Itoa(i),
				"sort":      "forks",
				"direction": "des",
			}).
			Get(url)

		base.Must(err)
		_ = resp
		//fmt.Println(resp)
		//fmt.Printf("%+v\n", result)
		//fmt.Println(result)
		link_re := regexp.MustCompile(`https://api.github.com/users/(.*)$`)

		for _, repo := range followers {
			fans_username := link_re.FindStringSubmatch(repo.Url)
			if len(fans_username) > 0 {
				followers_list = append(followers_list, fans_username[1])
			}
		}
	}
	followers_list = base.RemoveDuplicated(followers_list)
	return followers_list
}

func GetFollowers(username string) {
	ch := make(chan *models.Userinfo)
	go GetUserInfo(username, ch)
	//fmt.Printf("共有%d个粉丝\n", userinfo.Followers)
	userinfo := <-ch
	followers_list := getFollower(username, userinfo.Followers)
	fmt.Println(followers_list)
	fmt.Println(len(followers_list))
}

func GetHotFans(username string, topN int) map[string]int {
	//获取关注者中的大v
	base_ch := make(chan *models.Userinfo)
	go GetUserInfo(username, base_ch)
	userinfo := <-base_ch

	fmt.Printf("共有%d个粉丝\n", userinfo.Followers)

	followers_list := getFollower(username, userinfo.Followers)
	follower_map := make(map[string]int)
	wg := &sync.WaitGroup{}
	ch := make(chan *models.Userinfo, 20)
	for _, user := range followers_list {
		//follower := GetUserInfo(user)
		wg.Add(1)
		go func(username string) {
			GetUserInfo(username, ch)
			defer wg.Done()
		}(user)

	}

	go func(ch chan *models.Userinfo) {
		wg.Wait()
		close(ch)
	}(ch)

	for follower := range ch {
		//fmt.Println(follower)
		if follower == nil {
			continue
		}
		if follower.HtmlUrl == "" {
			continue
		}
		follower_map[follower.HtmlUrl] = follower.Followers
	}
	follower_map_list := base.Sort(follower_map)

	result := make(map[string]int)
	for i := 0; i < topN; i++ {
		name := base.GetMapKey(follower_map_list[i])
		value := base.GetMapValue(follower_map_list[i])
		result[name] = value
		fmt.Printf("User: %s, its fans: %d\n", name, value)
	}

	return result

}
