package main

import(
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"math/rand"
	"io/ioutil"
	"encoding/base64"
	"golang.org/x/net/websocket"
	"github.com/garyburd/redigo/redis"
)
//数据库结构 key:Username value:Password Email PhoneNumber SaveTime PhoneIP ComputerIP
//HSET testName Password "testPasswd" Email "testEmail" PhoneNumber "testPhoneNumber" SaveTime "testSaveTime" PhoneIP "testPhoneIP" ComputerIP "testComputerIP"
var RedisClient redis.Conn
func main(){
	var err error
	RedisClient, err = redis.Dial("tcp", "127.0.0.1:5520")
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
	defer RedisClient.Close()
	//RedisClient.Do("flushall")
	http.Handle("/Work",websocket.Handler(RequestWebSocket))
	http.HandleFunc("/Login",Login)
	http.HandleFunc("/Register",Register)
	http.HandleFunc("/UserInfo",UserInfo)
	http.HandleFunc("/GetData",GetData)
	http.HandleFunc("/SetData",SetData)
	http.HandleFunc("/test",test)
	if err := http.ListenAndServe(":6888", nil); err != nil {
		fmt.Println("ListenAndServe:", err)
    }
}
func GetData(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	UserID,found := r.Form["UserID"]
	if !found{
		fmt.Println("Get Error")
		w.Write([]byte("NoUserID"))
		return
	}
	RealUserIDBase64,err := base64.StdEncoding.DecodeString(UserID[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealUserID := string(RealUserIDBase64)
	fmt.Println("Get "+RealUserID + " New@"+RedisHMGET(RealUserID,"Data"))
	w.Write([]byte("New@"+RedisHMGET(RealUserID,"Data")))
}
func SetData(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	UserID,found := r.Form["UserID"]
	if !found{
		w.Write([]byte("NoUserID"))
		return
	}
	Data,found1 := r.Form["Data"]
	if !found1{
		w.Write([]byte("NoData"))
		return
	}
	RealUserIDBase64,err := base64.StdEncoding.DecodeString(UserID[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealUserID := string(RealUserIDBase64)
	RealDataBase64,err := base64.StdEncoding.DecodeString(Data[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealData := string(RealDataBase64)
	fmt.Println("Set "+RealUserID +" Data "+RealData)
	_,err1 := RedisClient.Do("HMSET",RealUserID,"Data",RealData)
    if err1 != nil {
		fmt.Println("redis hset error:", err)
		w.Write([]byte("SetError"))
	} else {
		//_,err := RedisClient.Do("expire","myKey","10")
		w.Write([]byte("SetSuccess"))
	}
}
func RequestWebSocket(ws *websocket.Conn){
	for{
		msg := make([]byte,512)
		n,err := ws.Read(msg)
		if err != nil{
			fmt.Println(err)
			break
		}
		ReceiveByte,_ := base64.StdEncoding.DecodeString(string(msg[:n]))
		ReceiveString := string(ReceiveByte)
		ReceiveData := strings.Split(ReceiveString,"|@|")
		if strings.Compare(ReceiveData[0],"Hi") == 0{
			SentMessages(ws,"Hai")
		} else if strings.Compare(ReceiveData[0],"AnNewText") == 0{
			if RedisSetValue(ReceiveData[1],ReceiveData[2]){
				SentMessages(ws,"AnNewTextOK")
			} else {
				SentMessages(ws,"AnNewTextFailed")
			}
		} else if strings.Compare(ReceiveData[0],"PCNewText") == 0{
			if RedisSetValue(ReceiveData[1],ReceiveData[2]){
				SentMessages(ws,"AnNewTextOK")
			} else {
				SentMessages(ws,"AnNewTextFailed")
			}
		}
	}
}
func SentMessages(ws *websocket.Conn,data string){
	SentMessage := base64.StdEncoding.EncodeToString([]byte(data))
	_, err := ws.Write([]byte(SentMessage))
	if err != nil {
		fmt.Println(err)
	}
}
func UserInfo(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	Username,found := r.Form["Username"]
	if !found{
		w.Write([]byte("NoUsername"))
		return
	}
	_,found1 := r.Form["Password"]
	if !found1{
		w.Write([]byte("NoPassword"))
		return
	}
	RealUsernameBase64,err := base64.StdEncoding.DecodeString(Username[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealUsername := string(RealUsernameBase64)
	Type,found2 := r.Form["Type"]
	if !found2{
		w.Write([]byte("error"))
		return 
	}
	is_key_exit,_ := redis.Bool(RedisClient.Do("EXISTS",RealUsername))
	if !is_key_exit{
		w.Write([]byte("NoUser"))
		return
	}
	if strings.Compare(Type[0],"Get") == 0{
		result, err := redis.Values(RedisClient.Do("HGETALL", RealUsername))
		ResultString := ""
		if err != nil {
			fmt.Println("hgetall failed", err.Error())
		} else {
			i := -1
			FirstEnter := true
			for _, v := range result {
				if i < 0{
					i = i * -1
				} else {
					if FirstEnter{
						ResultString = string(v.([]byte))
						FirstEnter = false
					} else {
						ResultString = ResultString +"|@|"+string(v.([]byte))
					}
					i = i * -1
				}
			}
			w.Write([]byte(ResultString))
			return
		}
	} else if strings.Compare(Type[0],"Set") == 0{
		Key,found3 := r.Form["Key"]
		if !found3{
			w.Write([]byte("error"))
			return 
		}
		Value,found4 := r.Form["Value"]
		if !found4{
			w.Write([]byte("error"))
			return 
		}
		_, err := RedisClient.Do("HMSET",RealUsername,Key[0],Value[0])
		if err!=nil{
			w.Write([]byte("UpdateError"))
			return
		}
		w.Write([]byte("UpdateSuccess"))
	}
}
func Login(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	Username,found := r.Form["Username"]
	if !found{
		w.Write([]byte("NoUserName"))
		return
	}
	Password,found1 := r.Form["Password"]
	if !found1{
		w.Write([]byte("NoPassword"))
		return
	}
	RealUsernameBase64,err := base64.StdEncoding.DecodeString(Username[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealPasswordBase64,err := base64.StdEncoding.DecodeString(Password[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealUsername := string(RealUsernameBase64)
	RealPassword := string(RealPasswordBase64)
	if RedisLoginCheck(RealUsername,RealPassword){
		w.Write([]byte("LoginSuccess@"+RedisHMGET(RealUsername,"UserID")))
	} else {
		w.Write([]byte("LoginError"))
	}
}
func Register(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	Username,found := r.Form["Username"]
	if !found{
		w.Write([]byte("NoUsername"))
		return
	}
	Password,found1 := r.Form["Password"]
	if !found1{
		w.Write([]byte("NoPassword"))
		return
	}
	var UserID string
	for i:=0;i<9;i++{
		UserID = UserID +strconv.Itoa(rand.Intn(10))
	}
	Email := "Please Set Your Email"
	PhoneNumber := "Please Set Your Pnone Number"
	SaveTime := "300"
	RealUsernameBase64,err := base64.StdEncoding.DecodeString(Username[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealPasswordBase64,err := base64.StdEncoding.DecodeString(Password[0])
	if err != nil{
		w.Write([]byte("error"))
		return 
	}
	RealUsername := string(RealUsernameBase64)
	RealPassword := string(RealPasswordBase64)
	is_key_exit,_ := redis.Bool(RedisClient.Do("EXISTS",RealUsername))
	if is_key_exit{
		w.Write([]byte("HaveUsername"))
		return
	} else {
		if RedisAddItem(RealUsername,UserID,RealPassword,Email,PhoneNumber,SaveTime) {
			w.Write([]byte("RegisterSuccess"))
		} else {
			w.Write([]byte("RegisterError"))
		}
	}
}
func test(w http.ResponseWriter, r *http.Request){
	b,_ := ioutil.ReadFile("test.html")
	w.Write(b)
}
func RedisSetValue(key string,value string)bool{
	_, err := RedisClient.Do("SET",key,value)
    if err != nil {
		fmt.Println("redis set failed:", err)
		return false
	}
	return true
}
func RedisGetValue(key string)string{
	Data,err := redis.String(RedisClient.Do("GET",key))
	if err != nil {
        fmt.Println("redis get failed:", err)
	}
	return Data
}
func RedisAddItem(Username string,UserID string,Password string,Email string,PhoneNumber string,SaveTime string)bool{
	_, err := RedisClient.Do("HMSET",Username,"UserID",UserID,"UserPassword",Password,"Email",Email,"PhoneNumber",PhoneNumber,"SaveTime",SaveTime)
    if err != nil {
		fmt.Println("redis hset error:", err)
		return false
	} else {
		//_,err := RedisClient.Do("expire","myKey","10")
		return true
	}
}
func RedisLoginCheck(UserName string,UserPassword string)bool{
	is_key_exit,_ := redis.Bool(RedisClient.Do("EXISTS",UserName))
	if !is_key_exit{
		return false;
	}
	PasswordResult,err := redis.Values(RedisClient.Do("HMGET",UserName,"UserPassword"))
	if err != nil{
		return false
	} else {
		var Password []byte
		for _, v := range PasswordResult {
			Password = v.([]byte)
			break
		}
		if strings.Compare(UserPassword,string(Password)) == 0{
			return true
		} else {
			return false
		}
	}
}
func RedisHMGET(key string,para string)string{
	Result,err := redis.Values(RedisClient.Do("HMGET",key,para))
	if err != nil{
		return "NULL"
	} else {
		var ResultString []byte
		for _, v := range Result {
			ResultString = v.([]byte)
			break
		}
		return string(ResultString)
	}
}
func RedisRemoveAllData()bool{
	_,err := RedisClient.Do("FLUSHALL")
	if err != nil{
		return false
	} else {
		return true
	}
}