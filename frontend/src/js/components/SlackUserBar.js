import React, {Component} from 'react';
import UserStore from '../userStore';
import PropTypes from 'prop-types';
import {Button, FormControl, InputGroup, Navbar} from 'react-bootstrap';

class SlackUserBar extends Component {
  constructor(props) {
    super(props);

    this.state = {
      editMode: false,
      email: ''
    };
  }

  saveUser() {
    if (this.isUserValid()) {
      const email = this.state.email;
      UserStore.saveUser(email);
      this.edit(false);
      this.props.onUserChange(email);
    } else {
      //TODO make it more user friendly
      alert('Please fill in valid email address.');
    }
  }

  edit(show) {
    this.setState({
      editMode: show
    });
  }

  toEdit() {
    if (UserStore.isUserSet()) {
      const user = UserStore.getUser();
      this.setState({
        email: user
      });
    }

    this.edit(true);
  }

  isUserValid() {
    return this.state.email !== '';
  }

  renderView() {
    const user = UserStore.getUser();
    return (
      <>
        <InputGroup>
          <FormControl disabled={true} type="email" value={user}/>
          <InputGroup.Append>
            <Button variant="light" onClick={() => this.toEdit()}>Edit</Button>
          </InputGroup.Append>
        </InputGroup>
      </>
    );
  }

  renderEdit() {
    const userSet = UserStore.isUserSet();
    return (
      <>
        <InputGroup>
          <FormControl type="email"
            placeholder="Slack user email"
            onChange={e => this.setState({email: e.target.value})}
            value={this.state.email}/>
          <InputGroup.Append>
            <Button variant="light" onClick={() => this.saveUser()}>Save</Button>
            {userSet && <Button variant="light" onClick={() => this.edit(false)}>Back</Button>}
          </InputGroup.Append>
        </InputGroup>
      </>
    );
  }

  navBar(comp) {
    return (
      <Navbar bg="dark" variant="dark">
        <Navbar.Brand>Happiness Door Bot</Navbar.Brand>
        <Navbar.Toggle/>
        <Navbar.Collapse className="justify-content-end">
          <Navbar.Text>
            {comp}
          </Navbar.Text>
        </Navbar.Collapse>
      </Navbar>
    );
  }

  render() {
    let userComponent;
    if (this.state.editMode || !UserStore.isUserSet()) {
      userComponent = this.renderEdit();
    } else {
      userComponent = this.renderView();
    }
    return this.navBar(userComponent);
  }
}

SlackUserBar.propTypes = {
  onUserChange: PropTypes.func
};

export default SlackUserBar;
