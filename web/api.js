function wake(host, user, pass, handle)
{

    host = encodeURIComponent(host);
    
    user = encodeURIComponent(user);
    
    pass = encodeURIComponent(pass);



    var data = {"host":host,"user":user,"pass":pass}
    var apiData = JSON.stringify(data);

    $.ajax(
    {
        url: "http://127.0.0.1:8081",
        type: "POST",
        headers: {'x-api':apiData},
        success: handle
    });
}


