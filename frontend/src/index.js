import ReactDOM from 'react-dom';
import React from 'react';
import './css/main.css';
import './css/bootstrap.min.css';
import MainView from './js/view/Main';

const wrapper = document.getElementById('container');
wrapper ? ReactDOM.render(<MainView />, wrapper) : false;