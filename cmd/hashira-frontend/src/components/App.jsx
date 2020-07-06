import React, { Component } from 'react';
import Cookies from 'js-cookie';

interface Props {}
interface State {
  name: String,
  isConnectedToGoogle: boolean,
  isConnectedToTwitter: boolean,
}

export default class App extends Component<Props, State> {
  constructor(props) {
    super(props)
    this.state = {
      name: "",
      isConnectedToGoogle:  false,
      isConnectedToTwitter: false,
    }
  }

  componentDidMount() {
    this.getMyInfo()
  }

  getMyInfo() {
    let headers = new Headers(
      {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': Cookies.get('Authorization'),
        'Origin': 'http://localhost:8080',
      })

    fetch('http://localhost:8080/auth/v1/me', {
      method: 'GET',
      credentials: 'include',
      headers: headers,
    }).then(response => {
      if (response.ok) {
        return response.json()
      } else {
        throw new Error()
      }
    }).then(body => {
      this.setState(prevState => {
        return {
          "Name": body.name,
          "isConnectedToGoogle": body.google_id !== "",
          "isConnectedToTwitter": body.twitter_id !== "",
        }
      })
    }).catch(error => {
      console.log(error)
    });
  }

  render() {
    return (
      <div>
        { this.state.Name ? "Hello, " + this.state.Name + " !" : "Let's login with:" }
        <p><a href="http://localhost:8080/auth/google">login by google</a> {this.state.isConnectedToGoogle ? 'Connected!' : ''}</p>
        <p><a href="http://localhost:8080/auth/twitter">login by twitter</a> {this.state.isConnectedToTwitter ? 'Connected!' : ''}</p>
      </div>
    )
  }
}
