# gomgo
Gomgo is a toy http json api server powered by Golang and MongoDB and ZeroMQ. It does not use magic or not magic framework, but just Idiomatic HTTP Middleware and some Best Practices for Golang, e.g. middleware chain, error handling, jwt token. I use it as a micro project template myself.


## Run server
1 Install Go and set up your [GOPATH](http://golang.org/doc/code.html#GOPATH), then install MongoDB and run it

2 Install Go package
~~~
go get github.com/gorilla/mux
go get github.com/justinas/alice
go get github.com/stretchr/graceful
go get github.com/PuerkitoBio/throttled
go get github.com/PuerkitoBio/throttled/store
go get gopkg.in/mgo.v2
go get gopkg.in/mgo.v2/bson
go get code.google.com/p/go.crypto/bcrypt
go get github.com/dgrijalva/jwt-go
~~~

3 Get this code
~~~
go get github.com/fengjian0106/gomgo
~~~

4 Build
~~~
sh build.sh
~~~

5 Run the server
~~~
./gomgo
~~~
You will now have a Go net/http webserver running on `localhost:3000`.



## Request the server
Use curl to make http request

1 Register new user
~~~
curl -v -X POST -H 'Content-Type: application/json' \
     -d '{"email": "helloworld@gmai.com", "password": "123456", "name": "helloworld"}' \
     http://127.0.0.1:8080/api/users
~~~

2 Login
~~~
curl -v -X POST -H 'Content-Type: application/json'  \
     -d '{"email": "helloworld@gmai.com", "password": "123456"}'  \
     http://127.0.0.1:8080/api/signin
~~~

After login, you will get response like below
~~~
{"userId":"53e49b07c3666ed09d000001","token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYWRkaW5nIjoie1wiaWRcIjpcIjUzZTQ5YjA3YzM2NjZlZDA5ZDAwMDAwMVwiLFwibmFtZVwiOlwiaGVsbG93b3JsZFwifSIsIlRva2VuVHlwZSI6IkFjY2Vzc1Rva2VuIiwiZXhwIjoxNDEwMDgyOTU0fQ.XN8OsZJKzv0HdloP-6T53PY4eOA8v59zBeIf_-6F0lRAoUqGpT6kgyuisaDlDDU_KkFiubOS2Akg0lj_sls7XkJJCR5sDgCHV9pRAhK41c9OEvq1OmJl0uxbOh22WOtbTLtyi_H6rS5Rxe3lOiL7dS539uLgBTzQshnXxXEWnQVKTFbJB2DitVnZNuAZTEKxjp1sbXBsLWDQ3IdfVwHRY8gX2g5f44QMBx83Qd-yvf0kIv-_bBugX7LXzruihKI8-caUsuaDAi--MoAqmVVsTHCImJvLjyZIhMqaRZSry48qo4NPCgUqoZOSQ9QkxQ0N1jWuGL9ahAL5Wgr5qzwv9g"}
~~~

Copy `userId` and `token`, later we will use it

3 Create a blog post  
Use `userId` and `token` you copied in last step 
~~~
curl -v -X POST -H 'Content-Type: application/json'  \
     -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYWRkaW5nIjoie1wiaWRcIjpcIjUzZTQ5YjA3YzM2NjZlZDA5ZDAwMDAwMVwiLFwibmFtZVwiOlwiaGVsbG93b3JsZFwifSIsIlRva2VuVHlwZSI6IkFjY2Vzc1Rva2VuIiwiZXhwIjoxNDEwMDgyOTU0fQ.XN8OsZJKzv0HdloP-6T53PY4eOA8v59zBeIf_-6F0lRAoUqGpT6kgyuisaDlDDU_KkFiubOS2Akg0lj_sls7XkJJCR5sDgCHV9pRAhK41c9OEvq1OmJl0uxbOh22WOtbTLtyi_H6rS5Rxe3lOiL7dS539uLgBTzQshnXxXEWnQVKTFbJB2DitVnZNuAZTEKxjp1sbXBsLWDQ3IdfVwHRY8gX2g5f44QMBx83Qd-yvf0kIv-_bBugX7LXzruihKI8-caUsuaDAi--MoAqmVVsTHCImJvLjyZIhMqaRZSry48qo4NPCgUqoZOSQ9QkxQ0N1jWuGL9ahAL5Wgr5qzwv9g" \
     -d '{"from": {"id": "53e49b07c3666ed09d000001"}, "message": "this is the first blog, hello world"}'  \
     http://127.0.0.1:8080/api/users/53e49b07c3666ed09d000001/posts
~~~

If you post blog success, you will get response like below
~~~
{"postId": "53e49e68c3666ed09d000002"}
~~~

Copy `postId`, later we will use it

4 Create a comment for a blog  
Normal, somebody else will create a comment for your post. Inorder to demonstrate this, you need register another user, and repeat step <1> and <2>, then you will get a new pair of `userId` and `token`. Copy them, also copy the `postId` you get in step <3>, we will use them 
~~~
curl -v -X POST -H 'Content-Type: application/json'  \
     -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYWRkaW5nIjoie1wiaWRcIjpcIjUzZTRhMWU1YzM2NjZlZDA5ZDAwMDAwM1wiLFwibmFtZVwiOlwibmV3b25lXCJ9IiwiVG9rZW5UeXBlIjoiQWNjZXNzVG9rZW4iLCJleHAiOjE0MTAwODQ2NjZ9.lKGbdz6zYu5aXpk2Xq5JLYfcBG4kpD5Wa7NiExSQYddtawX_5rosAFsFzD-mNbmR79Ymtt8j22kuw-mC1vZbCx6BCgMtqtb9X3Q7pvEKZ-46WKSJ5E2sHGNZ3YZa-iTYIf4CD0_LmWeHT5UPm3MWYo14Hf-tr6sLUeovmp7NuXj0x-pJDogSJ815NctoWFHXVTcTwffd52WaPptQjeryisROo1qbtmjPAAgdXKFBDWiwe2nrzG4erpbxOiGAOy9CT5rUhMiqlCKC-FGhc4UZ9GQ6pnzbv72-5uQqfiEJc3EWSuSbuyrNa-CAHDapr90SN3j3hLrE45PNVpQxotubFg" \
     -d '{"from":{"id": "53e4a1e5c3666ed09d000003"}, "message": "this is comment from other user"}'  \
     http://127.0.0.1:8080/api/posts/53e49e68c3666ed09d000002/comments
~~~

If you post comment success, you will get response like below
~~~
{"postId": "53e49e68c3666ed09d000002", "commentId": "53e4a25bc3666ed09d000004"}
~~~

5 Get the post  
Everyone can get post by the postId, so let's use the new `token` received in step <4> to get the post with the `postId` received in step <3>
~~~
curl -v --compressed -X GET -H 'Content-Type: application/json'  \
     -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYWRkaW5nIjoie1wiaWRcIjpcIjUzZTRhMWU1YzM2NjZlZDA5ZDAwMDAwM1wiLFwibmFtZVwiOlwibmV3b25lXCJ9IiwiVG9rZW5UeXBlIjoiQWNjZXNzVG9rZW4iLCJleHAiOjE0MTAwODQ2NjZ9.lKGbdz6zYu5aXpk2Xq5JLYfcBG4kpD5Wa7NiExSQYddtawX_5rosAFsFzD-mNbmR79Ymtt8j22kuw-mC1vZbCx6BCgMtqtb9X3Q7pvEKZ-46WKSJ5E2sHGNZ3YZa-iTYIf4CD0_LmWeHT5UPm3MWYo14Hf-tr6sLUeovmp7NuXj0x-pJDogSJ815NctoWFHXVTcTwffd52WaPptQjeryisROo1qbtmjPAAgdXKFBDWiwe2nrzG4erpbxOiGAOy9CT5rUhMiqlCKC-FGhc4UZ9GQ6pnzbv72-5uQqfiEJc3EWSuSbuyrNa-CAHDapr90SN3j3hLrE45PNVpQxotubFg" \
     http://127.0.0.1:8080/api/posts/53e49e68c3666ed09d000002
~~~

If everything is ok, you will get response like below
~~~
{"id":"53e49e68c3666ed09d000002","from":{"id":"53e49b07c3666ed09d000001","name":"helloworld"},"Message":"this is the first blog, hello world","CreatedTime":"2014-08-08T17:54:48.317+08:00","UpdatedTime":"2014-08-08T17:54:48.317+08:00","Comments":[{"id":"53e4a25bc3666ed09d000004","From":{"id":"53e4a1e5c3666ed09d000003","name":"newone"},"Message":"this is comment from other user","CreatedTime":"2014-08-08T18:11:39.752+08:00"}]}
~~~

## Distributed Systems
Distributed Systems is interesting. I will also try to show some basic technique on how to implement it.

1 RPC  
Path "/api/search?q=xxx&timeout=1s" do a google search, and the handler call a remote REST service (If you want a RPC, the code will be same the like)

You can test it like below, using keyword "golang"
~~~
curl -v --compressed -X GET -H 'Content-Type: application/json' http://127.0.0.1:8080/api/search?q=golang&timeout=2s
~~~

And if success, you will get response like below
~~~
{"data":[{"title":"The Go Programming Language","url":"http://golang.org/"},{"title":"A Tour of Go","url":"http://tour.golang.org/"},{"title":"Downloads - The Go Programming Language","url":"http://golang.org/dl/"},{"title":"Go (programming language) - Wikipedia, the free encyclopedia","url":"http://en.wikipedia.org/wiki/Go_(programming_language)"}],"elapsedSeconds":3.29384076}
~~~

2 ZeroMQ  
Path "/api/zmp?msg=xxx" send a request and get a reply. And the architecture is shown as below  
~~~
+-----------------------------------+                                                                                                            
|                                   |                                                                                                            
|  +----------+       +----------+  |                                                                                                            
|  |ZMQHandler|  ...  |ZMQHandler|  |                                                                                                            
|  +----------+       +----------+  |                                                                                                            
|                                   |                                                                                                            
|       ^                  ^        |                                                                                                            
|       |       gomgo      |        |                                                                                                            
|       |                  |        |                                                                                                            
|       v                  v        |                                     +--------+                                                             
|                                   |                                +--> | worker |                                                             
|     +----------------------+      |                                |    +--------+                                                             
|     |     go msgQueue()    |  <-------+      +---------------+     |                                                                           
|     +----------------------+      |   |      |               | <---+                                                                           
|                                   |   +----> |               |                                                                     
+-----------------------------------+          |  reqRepBroker |            ......                                                          
                                        +----> |               |                                                                      
                                        |      |               | <---+                                                                           
               ......                   |      +---------------+     |                                                                           
                                        |                            |    +--------+                                                             
                                        |                            +--> | worker |                                                             
+-----------------------------------+   |                                 +--------+                                                             
|              client               | <-+                                                                                                        
+-----------------------------------+                                                                                      
~~~                       

##    

Enter the dir zmqReqRepBrokerServer/nodejs-server, and install node.js dependency  
~~~
npm install
~~~

After run gomgo, you should also run ONLY ONE zeromq reqRepBroker  
~~~
node reqRepBroker.js
~~~

And run ONE or MANY zeromq worker 
~~~
node worker.js
~~~

And optional, you can run ONE or MANY zeromq client  
~~~
node client.js
~~~

Finally, it's time to test it with curl
~~~
curl -v --compressed -X GET -H 'Content-Type: application/json' http://127.0.0.1:8080/api/zmq?msg=helloworld
~~~

And if success, you will get response like below
~~~
{"retMsg": "node.js server [81043] echo: helloworld [1]"}
~~~

## Thank these guys and their articles!
[http://www.alexedwards.net/blog/a-recap-of-request-handling](http://www.alexedwards.net/blog/a-recap-of-request-handling)

[http://capotej.com/blog/2013/10/07/golang-http-handlers-as-middleware/](http://capotej.com/blog/2013/10/07/golang-http-handlers-as-middleware/)

[http://justinas.org/writing-http-middleware-in-go/](http://justinas.org/writing-http-middleware-in-go/)

[http://justinas.org/alice-painless-middleware-chaining-for-go/](http://justinas.org/alice-painless-middleware-chaining-for-go/)

[http://elithrar.github.io/article/custom-handlers-avoiding-globals/](http://elithrar.github.io/article/custom-handlers-avoiding-globals/)

[Go Concurrency Patterns: Context](http://blog.golang.org/context)

[http://angular-tips.com/blog/2014/05/json-web-tokens-introduction/](http://angular-tips.com/blog/2014/05/json-web-tokens-introduction/)

[http://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api](http://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api)

[https://github.com/soygul/koan](https://github.com/soygul/koan)


## License
MIT
