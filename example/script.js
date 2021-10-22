import Coral from "./Coral.js";

let Username = "";

//coral communication instance
const coral = new Coral()

//on connection
coral.onOpen = function () {
    // on message type
    coral.on("forward", (msg) => {
        addExternalMessage(msg.username, msg.message)
    })
}

// emit message to coral server
function sendInput(msg) {
    coral.emit("message", {username: Username, message: msg});
    addMyMessage(Username, msg);
}


function scrollToBottom() {
    const el = document.getElementById('messages')
    el.scrollTop = el.scrollHeight
}

setUsernameVal("user " + Date.now());


document.querySelector("#user-read").onclick = editUsername
document.querySelector("#user-edition input").onkeypress = saveUsername

function editUsername() {
    document.querySelector("#user-edition").classList.remove("hidden");
    document.querySelector("#user-edition input").focus()
    document.querySelector("#user-edition input").select()
    this.classList.add("hidden");

}

function saveUsername() {
    if (event.key == "Enter") {
        const u = getUsernameVal()
        if (u.length == 0) {
            return
        }
        document.querySelector("#user-edition").classList.add("hidden");
        document.querySelector("#user-read").classList.remove("hidden");
        setUsernameVal(u)
    }

}

function getUsernameVal() {
    return document.querySelector("#username-rw").value;

}

function setUsernameVal(str) {
    document.querySelector("#username-rw").value = str;
    document.querySelector("#username-ro").innerHTML = str;
    Username = str
}

document.querySelector("#text-input").onkeypress = (event, element) => {
    if (event.key == "Enter" && !event.shiftKey) {
        const msg = htmlEntities(event.target.value)
        if (msg.trim().length > 0) {
            sendInput(msg);
        }
        event.target.value = ""


    }
}
document.querySelector("#text-input").onkeyup = (event, element) => {
    if (event.key == "Enter" && !event.shiftKey) {
        event.target.value = ""
    }
}



function addMyMessage(username, content) {
    content = nl2br(content)
    let div = document.createElement('div');

    div.innerHTML = `
                   <div class="flex items-end justify-end">
                    <div class="flex flex-col space-y-2 text-xs max-w-xs mx-2 order-1 items-end">
                        <div><div class="px-4 py-2 rounded-lg inline-block bg-red-300 text-white ">

                            <div>${content}
                            </div>
                            </div>
                        </div>
                    </div>
                </div>`;
    document.querySelector("#messages").appendChild(div);
    scrollToBottom();

}

function addExternalMessage(username, content) {
    content = nl2br(content)
    let div = document.createElement('div');

    div.innerHTML = `<div class="flex items-end">
                    <div class="flex flex-col space-y-2 text-xs max-w-xs mx-2 order-2 items-start">
                        <div>
                            <div class="flex flex-col px-4 py-2 rounded-lg inline-block bg-gray-300 text-gray-600">
                            <span class="font-black">${username}</span>

                            <div>${content}
                            </div>
                            </div>
                        </div>
                    </div>
                </div>`;
    document.querySelector("#messages").appendChild(div);
    scrollToBottom();
}

function nl2br(str) {
    return (str).replace(/([^>\r\n]?)(\r\n|\n\r|\r|\n)/g, '$1<br/>$2');
}

function htmlEntities(str) {
    return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}