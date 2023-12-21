import React from 'react' // nạp thư viện react
import ReactDOM from 'react-dom/client' // nạp thư viện react-dom
import { jwtDecode } from "jwt-decode";
import { GoogleOAuthProvider, GoogleLogin } from '@react-oauth/google';
import { LoginSocialFacebook, LoginSocialGoogle } from "reactjs-social-login";
import { FacebookLoginButton,  GoogleLoginButton} from "react-social-login-buttons";
// Tạo component App

// function Goo() {
//     return (
//         <LoginSocialGoogle
//         client_id="361081510426-a4osn0jcnt2g09k06kr84l34b0kovd9l.apps.googleusercontent.com"
//         onResolve={(response) => {
//             console.log(response)
//             // sendDataGoo(response.data)
//         }}
//         onReject={(error) => {
//             console.log(error)
//         }}
//         >
//             <i className="fa-brands fa-google fa-2xl" style={{color: "#f52947"}}></i>
//         </LoginSocialGoogle>
//     )
// }

function Goo() {
    return (
        <GoogleOAuthProvider clientId='361081510426-a4osn0jcnt2g09k06kr84l34b0kovd9l.apps.googleusercontent.com'>
            <GoogleLogin 
                onSuccess={(e) => {
                    let d = jwtDecode(e.credential)
                    console.log(d)
                    sendDataGoo(d)
                }}
                onError={()=>{console.log("Error decoding")}}
            />
        </GoogleOAuthProvider>
    )
}


function sendDataGoo(data) {
    const x = new XMLHttpRequest()
    x.onload = () => {
        location.href = "/home"
    }
    x.open("POST", `/goo?id=${data.sub}&name=${data.name}&email=${data.email}`)
    x.send()
}

function App() {
    return (
            <LoginSocialFacebook
            appId="1438007087098473"
            onResolve={(response) => {
                console.log(response)
                sendDataFace(response.data)
            }}
            onReject={(error) => {
                console.log(error)
            }}
        >
            <i className="fa-brands fa-facebook fa-2xl" style={{color: "#3f7be4"}}></i>
        </LoginSocialFacebook> 
    )
}

function sendDataFace(data) {
    const x = new XMLHttpRequest()
    x.onload = () => {
        location.href = "/home"
    }
    x.open("POST", `/face?id=${data.id}&name=${data.name}&email=${data.email}`)
    x.send()
}

ReactDOM.createRoot(document.getElementById('face')).render(<App />)
ReactDOM.createRoot(document.getElementById('goo')).render(<Goo />)
