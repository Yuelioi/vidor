package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	ChunkSize  int64
	client     clientInfo
}

var defaultClient = clientInfo{
	userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	referer:   "https://www.bilibili.com",
}

func (c *Client) assureClient() {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
}

type clientInfo struct {
	userAgent string
	referer   string
	sessdata  string
}

func NewClient(sessdata string) *Client {
	defaultClientCopy := defaultClient    // 创建一个 defaultClient 的副本
	defaultClientCopy.sessdata = sessdata // 修改副本的 sessdata 字段
	return &Client{
		HTTPClient: http.DefaultClient,
		client:     defaultClientCopy,
	}
}
func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("Get", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.client.userAgent)
	req.Header.Set("Referer", c.client.referer)
	req.Header.Set("Cookie", "SESSDATA="+c.client.sessdata)
	return req, nil

}

func (c *Client) Response(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil, err
	}
	return resp, nil
}

func (c *Client) Get(url string, body io.Reader) ([]byte, error) {
	resp, err := c.Response("Get", url, nil)

	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	return data, nil
}

func (c *Client) GetPlaylistInfo(aid int, bvid string) (*biliPlaylistInfo, error) {
	var bv biliPlaylistInfo

	playListApiUrl := fmt.Sprintf(`%s/x/web-interface/view?aid=%d&bvid=%s`, apiURL, aid, bvid)

	body, err := c.Get(playListApiUrl, nil)
	if err != nil {
		fmt.Println("Error cannot fetch Aid:", err)
		return nil, err
	}
	err = json.Unmarshal(body, &bv)
	if err != nil {
		fmt.Println("Error  JSON:", err)
		return nil, err
	}

	if bv.Code != 0 {
		return nil, errors.New(bv.Message)
	}

	return &bv, nil

}

func (c *Client) GetVideoDownloadInfo(bvid string, cid int) (*biliDownloadInfo, error) {

	downloadApiUrl := fmt.Sprintf("%s/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=4048&fourk=1", apiURL, bvid, cid)
	body, err := c.Get(downloadApiUrl, nil)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	bi := biliDownloadInfo{}
	if err = json.Unmarshal(body, &bi); err != nil {
		return nil, err
	}
	return &bi, nil
}
