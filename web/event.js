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
    else if(response == "fail")
    {
        $("#error").show();
        $("#outer").css({"top":"80px"});
        $("#error").html("host not found");        
    }
    return null;
}

$(document).ajaxError(serverDown);

function serverDown()
{
    $("#error").show();
    $("#outer").css({"top":"80px"});
    $("#error").html("Service Unavailable"); 
}
