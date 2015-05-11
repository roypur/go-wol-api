$(document).ready(
function()
{
    $("#btn").click(function()
    {
        wake($("#host").val(), $("#user").val(),$("#pass").val(),error);
        console.log(eval($("#host").val()));
    })
});
function error(response)
{
    if(response == "ok")
    {
        $("#error").hide();
        $("#outer").css({"top":"130px"});
    }
    else if(response == "denied")
    {
        $("#error").show();
        $("#outer").css({"top":"80px"});
        $("#error").html("Wrong username or password");
    }
    else if(response == "fail")
    {
        $("#error").show();
        $("#outer").css({"top":"80px"});
        $("#error").html("host not found");        
    }
    else if(response == "fail")
    {
        $("#error").show();
        $("#outer").css({"top":"80px"});
        $("#error").html("host not found");        
    }
    else if(response == "down")
    {
        $("#error").show();
        $("#outer").css({"top":"80px"});
        $("#error").html("Service Unavailable");     
    }
    return null;
}

function serverDown()
{
    $("#error").show();
    $("#outer").css({"top":"80px"});
    $("#error").html("Service Unavailable");
}



function wake(host, user, pass, handle)
{

    var data = {"host":host,"user":user,"pass":pass}
    var apiData = f.stringify(data);

    f.api(apiData, ["http://oxygen.purser.it:8781","http://home.purser.it:8781"], error, serverDown);

}
