/*
New user validation




*/
console.log("User validation is working")

// Will return a bool or null - testing whether the username is unique
async function UsernameExists(username) {
    try {
        const result = await fetch(`/api/users/${username}`)
        return result;
    } catch(error) {
        alert(error);
        return false;
    }
}

let exists;
UsernameExists("89benmatthews@gmail.com").then(data => {
    exists = data
    console.log(exists);
});

//document.querySelectorAll() // something that connects the 