var apiUrl = "";

$(document).ready(
    function(){
        
        apiUrl = $("api-url").html()    
    
        notLoggedIn()
    
        $("#wake-btn").click(function()
        {
            wake($("#host").val(), $("#user").val(),$("#pass").val(),error);
        })
        
        $("#login-btn").click(function()
        {
            list($("#user").val(),$("#pass").val(),error);
            
        })        
    }
);
function error(response){

    if(response == "ok"){
        $("#error").hide();
        $("#outer").css({"top":"130px"});
    }else if(response == "denied"){
        $("#error").show();
        $("#outer").css({"top":"80px"});
        $("#error").html("Wrong username or password");
        notLoggedIn()
        
    }else if(response == "fail"){
        $("#error").show();
        $("#outer").css({"top":"80px"});
        $("#error").html("host not found");        
    }else{
    
        $(".not-logged-in").hide();
        $(".logged-in").show();
        $("#outer").css({"height":"250px"});
        
    }

    return null;
}

function serverDown(){
    $("#error").show();
    $("#outer").css({"top":"80px"});
    $("#error").html("Service Unavailable");
}

function notLoggedIn(){
    $(".logged-in").hide()
    $(".not-logged-in").show()

    $("#outer").css({"height":"300px"});
    
}

function wake(host, user, pass, handle){

    var req = $.ajax({
        url: apiUrl,
        crossDomain: true,
        headers: {"Goapi-User": user, "Goapi-Pass": pass, "Goapi-Host": host},
        timeout: 1500
        
    })

    req.done(handle);
    req.fail(serverDown)
}
function list(user, pass, handle){


    var req = $.ajax({
        url: apiUrl,
        crossDomain: true,
        headers: {"Goapi-User": user, "Goapi-Pass": pass, "Goapi-List": "true"},
        timeout: 1500
        
    })
    
    req.done(function(r){
    
        var decoded = JSON.parse(r)

        var options
            
        for(var i=0; i<decoded.length; i++){
            options += "<option value=\"" + decoded[i] + "\">"
        }
        
        $("#hosts").html(options);

        handle(r)    
    });
    req.fail(serverDown)
}
