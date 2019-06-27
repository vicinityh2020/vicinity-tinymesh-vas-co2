import React from 'react';
import './css/App.css';
import LineChartCO2 from './Components/LineChartCO2';
import {BrowserRouter, Route, Switch} from 'react-router-dom';
import {Col, Container, Row} from 'react-bootstrap';

function Home() {
    return (
        <Container>
            <h2>TinyMesh AS - VAS CO2</h2>
            <Row>
                <Col style={{paddingBottom: '30px'}} sm={12}>
                    <LineChartCO2 sensor={
                        {name: 'LAS00016222', oid: 'aac3fff0-49d6-45dd-aa3f-77c3e36644c8'}
                    }/>
                </Col>
                <Col sm={12}>
                    <LineChartCO2 sensor={
                        {name: 'LAS00016225', oid: '6d7f79e5-f8f5-4bfb-b8b4-cd09fea86bbb'}
                    }/>
                </Col>
            </Row>
        </Container>
    );
}

function App() {
    return (
        <BrowserRouter>
            <Switch>
                <Route path="/" exact component={Home}/>
            </Switch>
        </BrowserRouter>
    );
}

export default App;
