// Copyright 2016 Google, Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
$(document).ready(function () {


    Handlebars.registerHelper('each', function(context, options) {
        var ret = "";

        console.log(context)
        console.log(context.length)

        for(var i=0, j=context.length; i<j; i++) {
            console.log(context[i])
            ret = ret + options.fn(context[i]);
        }

        return ret;
    });


    var userstemplate = document.getElementById('userstemplate').innerHTML;
    var compiled_template = Handlebars.compile(userstemplate);
    var rendered = compiled_template({Users: []});
    document.getElementById('userstable').innerHTML = rendered;


    if ('WebSocket' in window) {
        var link = document.createElement("a");
        link.href = "/";
        var ws_host = "ws://"+link.host+"/ws"
        connect(ws_host);
    }
    else {
        alert("WebSockets don't seem to be supported on this browser.");
    }

    function sendData(ws) {
        return function () {
            textarea = document.getElementById('textdata')
            console.log("Sending")
            console.log({"Text": textarea.value})
            ws.send({"Text": textarea.value})
        };
    }

    function connect(host) {
        ws = new WebSocket(host);

        ws.onopen = function () {
            console.log('Connected!');
            textarea = document.getElementById('textdata')
            textarea.addEventListener("keyup", sendData(ws));
            textarea.addEventListener("paste", sendData(ws));
        };

        ws.onmessage = function (evt) {
            data = JSON.parse(evt.data)
            console.log(data);
            var rendered = compiled_template(data);
            document.getElementById('userstable').innerHTML = rendered;

            //document.getElementById('textdata').value = data.Text;
            textarea = document.getElementById('textdata')
            if (textarea == document.activeElement) {
                console.log("textarea is focused");
            } else {
                textarea.value = data.Text;
            }


        };

        ws.onclose = function () {
            console.log('Socket connection was closed!!!');
            connect(ws_host);
        };
    };
});
