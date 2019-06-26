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
                        {name: 'LAS00016222', oid: '5b0a4b48-d71e-4cf6-9a73-022e2cedb7e1'}
                    }/>
                </Col>
                <Col sm={12}>
                    <LineChartCO2 sensor={
                        {name: 'LAS00016225', oid: '57543cd0-5215-4667-b89d-24968a503c6b'}
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
