import React, {useEffect, useState} from 'react';
import './css/App.css';
import LineChartCO2 from './Components/LineChartCO2';
import {BrowserRouter, Route, Switch} from 'react-router-dom';
import {Col, Container, Row} from 'react-bootstrap';
import axios from 'axios';

function Home() {
    const [sensors, setSensors] = useState([]);

    useEffect(() => {
        const fetchSensors = async () => {
            const options = {
                method: 'get',
                url: '/api/objects',
                time: 5000
            };

            return await axios(options);
        };

        fetchSensors().then((response) => {
            setSensors(response.data.sensors);
        }).catch((error) => {
            console.log(error);
        });
    }, []);

    const getSensors = () => {
        return sensors.map((e) => {
            return (
                <Col key={e.name} style={{paddingBottom: '30px'}} sm={12}>
                    <LineChartCO2 sensor={
                        {name: e.name, oid: e.oid}
                    }/>
                </Col>
            );
        });
    };

    return (
        <Container>
            <h2>TinyMesh AS - VAS CO2</h2>
            <Row>
                {getSensors()}
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
