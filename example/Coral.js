function getDefaultURL() {
    if (location.protocol == "https:") {
        return 'wss://' + location.host + "/ws"
    }
    return 'ws://' + location.host + "/ws"
}

export default class {
    /**
     *
     * @param params
     */
    constructor(params = {
        url: getDefaultURL(),
    }) {
        this.opened = false
        this.ws = new WebSocket(params.url);
        this.ws.onopen = () => {
            this.opened = true
            if (this.onOpen) {
                this.onOpen()
            }
            this.ws.onmessage = (message) => {
                const msg = JSON.parse(message.data)
                if (this.registeredActions[msg.type]) {
                    this.registeredActions[msg.type](msg.content)
                }
            }

        }
        this.ws.onclose = () => {
            if (this.onClose) {
                this.onClose()
            }
            this.opened = false
        }
        this.registeredActions = {}
    }

    on(type, callback = () => {
    }) {
        this.registeredActions[type] = callback
    }

    emit(type, content = null) {
        if (!this.opened) {
            throw "cannot emit message, connection closed"
        }
        this.ws.send(JSON.stringify({type, content}))
        return true
    }

    join(channelName) {
        this.emit("_coral.channel.join", {channel: channelName})
        let chan = new Channel(channelName, this.ws)
        return chan
    }

    leave(channelName) {
        this.emit("_coral.channel.leave", {channel: channelName})
    }


}

class Channel {

    onLeft(callback){
        this.leftAction=callback
        return this
    }

    onJoin(callback){
        this.joinedAction=callback
        return this
    }
    constructor(name = "", ws) {
        this.ws = ws
        this.name = name
    }

    on(type, callback = () => {
    }) {
        this.registeredActions[type] = callback
    }

    emit(type, content = null) {
        if (!this.joined) {
            throw "cannot emit message, connection closed"
        }
        this.ws.send(JSON.stringify({type, content}))
        return true
    }

}