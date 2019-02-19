function fetchStream(dateStr, callback) {
    var url = "/api/list/"
    if (dateStr == null) {
        dateStr = (new Date()).toJSON()
    }

    url = url + dateStr

    fetch(url, {
        method: "GET",
        credentials: "include"
        })
        .then( response => response.json() )
        .then( response => {
            callback(response.Posts, response.EndOfList)
        })
}

function fetchPost(postId, callback) {
    var url = "/api/read/" + postId

    fetch(url, {
        method: "GET",
        credentials: "include"
        })
        .then( response => response.json() )
        .then( post => {
            callback(post)
        })
}

function createPost(text, callback) {
    var url = "/api/write"
    let post = {
        "Title": text
    }
    let postRequest = {
        "Content": JSON.stringify(post)
    }
    let body = JSON.stringify(postRequest)

    fetch(url, {
        method: "POST",
        credentials: "include",
        body: body, 
        credentials: "include"
        })
        .then( response => {
            if (!response.ok) {
                throw Error("Bad reply")
            }
            return response
        })
        .then( response => response.json() )
        .then( post => {
            callback(null, post)
        })
        .catch( error => {
            callback(error, null)
        })
}

function divForPost(post, inList) {
    let div = document.createElement("div")
    div.className = `post-box-${post.PostId} post-box`
    let a_open = inList ? "<a href='/r/" + post.PostId + "'>" : ""
    let a_close = inList ? "</a>" : ""
    div.innerHTML = "<div class='post'><div>" + a_open + post.Created + a_close + "</div><div>" + post.Title + "</div></div>"

    return div
}

function ShowError(error) {
    let errorBoxText = document.getElementById("error-box-text")
    errorBoxText.innerHTML = error
    errorBoxText.style = "display:block;"
}

function HideError() {
    let errorBoxText = document.getElementById("error-box-text")
    errorBoxText.style = "display:none;"
}
