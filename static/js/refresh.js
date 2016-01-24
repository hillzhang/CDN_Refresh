/**
 * Created by Hill on 2016/1/21.
 */
function refresh() {
    $.post('/refresh', {
        'url' : $("#url").val()
    }, function(json) {
        if (json.msg.length > 0) {
            err_message_quietly(json.msg);
        } else {
            ok_message_quietly('Refresh file successfully');
        }
    });
}

function err_message_quietly(msg, f) {
    $.layer({
        title : false,
        closeBtn : false,
        time : 2,
        dialog : {
            msg : msg
        },
        end : f
    });
}

function ok_message_quietly(msg, f) {
    $.layer({
        title : false,
        closeBtn : false,
        time : 2,
        dialog : {
            msg : msg,
            type : 1
        },
        end : f
    });
}