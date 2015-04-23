function wake(host, user, pass, handle)
{

    host = goEncode(host);
    
    user = goEncode(user);
    
    pass = goEncode(pass);

    console.log(pass);

    var data = {"host":host,"user":user,"pass":pass}
    var apiData = JSON.stringify(data);

    $.ajax(
    {
        url: "http://home.purser.it:8781",
        type: "POST",
        headers: {'x-api':apiData},
        success: handle
    });
}

function goEncode(str)
{
    str = str.replace('-', '-a');    
    str = str.replace('[', '-b');
    str = str.replace(']', '-c');
    str = str.replace('{', '-d');
    str = str.replace('}', '-e');
    str = str.replace("'", '-f');
    str = str.replace('"', '-g');
    str = str.replace(' ', '-h');
    str = str.replace(',', '-i');
    str = str.replace(':', '-j');
    
    return str
}

