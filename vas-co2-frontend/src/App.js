import React from 'react';
import './css/App.css';
import LineChartCO2 from './Components/LineChartCO2';
import {BrowserRouter, Switch, Route} from 'react-router-dom';
import {Col, Container, Row} from 'react-bootstrap';

function Home() {
    return (
        <Container>
            <Row>
                <Col sm={12}>
                    <LineChartCO2 props={
                        [
                            {uniqueId: 'LAS00016222', sensorOid: '5b0a4b48-d71e-4cf6-9a73-022e2cedb7e1'}
                            // {uniqueId: 'LAS00016225', sensorOid: '57543cd0-5215-4667-b89d-24968a503c6b'}
                        ]
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
                <Route path="/" exact component={Home} />
            </Switch>
        </BrowserRouter>
    );
}

export default App;
