function wake(host, user, pass, handle)
{

    host = f.encode(host);
    
    user = f.encode(user);
    
    pass = f.encode(pass);

    var data = {"host":host,"user":user,"pass":pass}
    var apiData = JSON.stringify(data);

    f.api(apiData, "http://home.purser.it:8781", handle);
}
