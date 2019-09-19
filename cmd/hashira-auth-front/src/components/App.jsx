import React, { Component } from 'react';
import Cookies from 'js-cookie';

interface Props {}
interface State {
  isConnectedToGoogle: boolean,
  isConnectedToTwitter: boolean,
}

export default class App extends Component<Props, State> {
  constructor(props) {
    super(props)
    this.state = {
      isConnectedToGoogle:  Cookies.get('GoogleID')  ? true : false,
      isConnectedToTwitter: Cookies.get('TwitterID') ? true : false,
    }
  }

  render() {
    return (
      <div>
        <p><a href="http://localhost:8080/auth/google">login by google</a> {this.state.isConnectedToGoogle ? 'Connected!' : ''}</p>
        <p><a href="http://localhost:8080/auth/twitter">login by twitter</a> {this.state.isConnectedToTwitter ? 'Connected!' : ''}</p>
      </div>
    )
  }
}
