import React, {useEffect, useState} from 'react';
import {Line} from 'react-chartjs-2';
import {Button} from 'react-bootstrap';

import axios from 'axios';

const data = {
    labels: [],
    datasets: [
        {
            label: 'My First dataset',
            fill: true,
            lineTension: 0.1,
            backgroundColor: 'rgba(75,192,192,0.4)',
            borderColor: 'rgba(75,192,192,1)',
            borderCapStyle: 'butt',
            borderDash: [],
            borderDashOffset: 0.0,
            borderJoinStyle: 'miter',
            pointBorderColor: 'rgba(75,192,192,1)',
            pointBackgroundColor: '#fff',
            pointBorderWidth: 1,
            pointHoverRadius: 5,
            pointHoverBackgroundColor: 'rgba(75,192,192,1)',
            pointHoverBorderColor: 'rgba(220,220,220,1)',
            pointHoverBorderWidth: 2,
            pointRadius: 1,
            pointHitRadius: 10,
            data: []
        }
    ]
};

function setYAxis(name) {
    return {
        scales: {
            yAxes: [{
                ticks: {
                    userCallback: function (item) {
                        return `${item} ${name}`;
                    },
                }
            }]
        }
    };
}

export default function LineChartCO2(sensors) {

    const [dataCO2, setDataCO2] = useState(data);
    const [toggle, setToggle] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            const options = {
                method: 'get',
                url: '/api/objects/57543cd0-5215-4667-b89d-24968a503c6b',
                time: 3000
            };

            const response = await axios(options);
            console.log(response);

            setDataCO2({
                labels: response.data.labels,
                datasets: [
                    {
                        label: 'My First dataset',
                        fill: true,
                        lineTension: 0.1,
                        backgroundColor: 'rgba(75,192,192,0.4)',
                        borderColor: 'rgba(75,192,192,1)',
                        borderCapStyle: 'butt',
                        borderDash: [],
                        borderDashOffset: 0.0,
                        borderJoinStyle: 'miter',
                        pointBorderColor: 'rgba(75,192,192,1)',
                        pointBackgroundColor: '#fff',
                        pointBorderWidth: 1,
                        pointHoverRadius: 5,
                        pointHoverBackgroundColor: 'rgba(75,192,192,1)',
                        pointHoverBorderColor: 'rgba(220,220,220,1)',
                        pointHoverBorderWidth: 2,
                        pointRadius: 1,
                        pointHitRadius: 10,
                        data: response.data.data
                    }
                ]
            });
        };

        fetchData().catch((error) => {
            console.log(error);
        });

    }, [toggle]);

    const toggleFetch = () => setToggle(!toggle);

    return (
        <div>
            <h2>Line Example</h2>
            <Button onClick={toggleFetch}/>
            <Line data={dataCO2} options={setYAxis('ppm')}/>
        </div>
    );
}