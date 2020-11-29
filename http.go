package z_cache

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zhaobing/z-cache/consistent_hash"
	pb "github.com/zhaobing/z-cache/zcachepb"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/zcache/"
	defaultRepicas  = 50
)

//HTTPPool 即作为http客户端发起请求，也作为http服务端接受请求
type HTTPPool struct {
	//ip:port  访问地址
	ipPort string
	//基础路由
	basePath string
	//guards peerPicker and httpGetters
	mu sync.Mutex
	//key与节点映射的哈希环
	peers *consistent_hash.HashCircle
	//远程节点与其对应httpGetter的映射,每个远程节点对应一个httpGetter
	httpGetters map[string]*httpGetter
}

//Set 实例化每一个节点
func (p *HTTPPool) SetPeers(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistent_hash.New(defaultRepicas, nil)
	p.peers.AddPhysicalNode(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// SelectPeer 根据具体的key,选择节点,返回节点对应的http客户端
func (p *HTTPPool) SelectPeer(key string) (peer PeerGetter, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	log.Println("SelectPeer,key=", key)
	if peer := p.peers.GetPhysicalNode(key); peer != "" && peer != p.ipPort {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	} else {
		fmt.Println("peer", peer)
	}
	return nil, false
}

//NewHTTPPool HTTPOOL构造
func NewHTTPPool(ipPortParam string) *HTTPPool {
	return &HTTPPool{
		ipPort:   ipPortParam,
		basePath: defaultBasePath,
	}
}

//Log 记录日志
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s ", p.ipPort, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}
	p.Log("Req:%s %s", r.Method, r.URL.Path)

	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := proto.Marshal(&pb.Response{
		Value: view.ByteSlice(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)

}

/** httpGetter ====================================== **/
type httpGetter struct {
	//表示将要访问的远程节点的地址，例如 http://example.com/zcache/
	baseURL string
}

//Get 访问远程节点,获取group和key对应的缓存
func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.Group),
		url.QueryEscape(in.Key),
	)
	log.Println("[getFromRemotePeer]", u)
	res, err := http.Get(u)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned:%v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body:%v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body:%v", err)
	}

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)
