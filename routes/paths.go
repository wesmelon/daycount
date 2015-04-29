package routes

import (
    "net/http"
)

var api = "/api/"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{"Index", "GET", "/", IndexPageHandler},
    Route{"Internal", "GET", "/internal", InternalPageHandler}, 

    //auth
    Route{"Login", "POST", "/login", LoginHandler}, 
    Route{"Logout", "POST", "/logout", LogoutHandler}, 
    Route{"Signup", "GET", "/signup", SignupHandler}, 
    Route{"GoSignup", "POST", "/signup", GoSignup}, 

    //api calls
    Route{"GetCategoriesByUser", "GET", api+"categories/user/{uid}", authHandler(GetCategoriesByUser)},
    Route{"GetCategoryById", "GET", api+"categories/{id}", authHandler(GetCategoryById)}, 
    Route{"PostCategory", "POST", api+"categories/", authHandler(PostCategory)}, 
    Route{"PutCategory", "PUT", api+"categories/{id}", authHandler(PutCategory)}, 
    Route{"DeleteCategory", "DELETE", api+"categories/{id}", authHandler(DeleteCategory)}, 

    Route{"GetContainersByUser", "GET", api+"containers/user/{uid}", authHandler(GetContainersByUser)},
    Route{"GetContainerById", "GET", api+"containers/{id}", authHandler(GetContainerById)}, 
    Route{"PostContainer", "POST", api+"containers/", authHandler(PostContainer)},
    Route{"PutContainer", "PUT", api+"containers/{id}", authHandler(PutContainer)},  
    Route{"DeleteContainer", "DELETE", api+"containers/{id}", authHandler(DeleteContainer)}, 

    Route{"GetDatesByContainer", "GET", api+"dates/container/{cid}", authHandler(GetDatesByContainer)},
    Route{"GetDateById", "GET", api+"dates/{id}", authHandler(GetDateById)}, 
    Route{"PostDate", "POST", api+"dates/", authHandler(PostDate)},
    Route{"PutDate", "PUT", api+"dates/{id}", authHandler(PutDate)},  
    Route{"DeleteDate", "DELETE", api+"dates/{id}", authHandler(DeleteDate)}, 
}