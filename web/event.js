$(document).ready(
function()
{
    $("#btn").click(function()
    {
        wake($("#host").val(), $("#user").val(),$("#pass").val(),error);
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

    f.api(apiData, ["http://home.royolav.net:2052"], error, serverDown);

}
