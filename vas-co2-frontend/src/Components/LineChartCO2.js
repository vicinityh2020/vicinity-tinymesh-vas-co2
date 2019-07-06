import React, {useEffect, useState} from 'react';
import {Chart, Line} from 'react-chartjs-2';

import {Button} from 'react-bootstrap';
import PropTypes from 'prop-types';

import axios from 'axios';

const MAX_X_LABELS = 11;

const data = {
    labels: [],
    datasets: [
        {
            label: 'No data for the last 12 hours',
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
                    padding: 0,
                    beginAtZero: true,
                }
            }],
            xAxes : [{
                ticks: {
                    autoSkip: true,
                    maxTicksLimit: MAX_X_LABELS
                }
            }]
        },
        animation: {
            duration: 0
        }
    };
}

LineChartCO2.propTypes = {
    sensor: PropTypes.object.isRequired,
};

function isEmpty(v) {
    return (typeof (v) === 'undefined' || v == null);
}

export default function LineChartCO2(props) {

    const [dataCO2, setDataCO2] = useState(data);
    const [toggle, setToggle] = useState(false);

    Chart.pluginService.register({
        afterUpdate: function (chart) {
            var xScale = chart.scales['x-axis-0'];
            if (xScale.options.ticks.maxTicksLimit) {
                // store the original maxTicksLimit
                xScale.options.ticks._maxTicksLimit = xScale.options.ticks.maxTicksLimit;
                // let chart.js draw the first and last label
                xScale.options.ticks.maxTicksLimit = (xScale.ticks.length % xScale.options.ticks._maxTicksLimit === 0) ? 1 : 2;

                var originalXScaleDraw = xScale.draw;
                xScale.draw = function () {
                    originalXScaleDraw.apply(this, arguments);

                    var xScale = chart.scales['x-axis-0'];
                    if (xScale.options.ticks.maxTicksLimit) {
                        var helpers = Chart.helpers;

                        var tickFontColor = helpers.getValueOrDefault(xScale.options.ticks.fontColor, Chart.defaults.global.defaultFontColor);
                        var tickFontSize = helpers.getValueOrDefault(xScale.options.ticks.fontSize, Chart.defaults.global.defaultFontSize);
                        var tickFontStyle = helpers.getValueOrDefault(xScale.options.ticks.fontStyle, Chart.defaults.global.defaultFontStyle);
                        var tickFontFamily = helpers.getValueOrDefault(xScale.options.ticks.fontFamily, Chart.defaults.global.defaultFontFamily);
                        var tickLabelFont = helpers.fontString(tickFontSize, tickFontStyle, tickFontFamily);
                        var tl = xScale.options.gridLines.tickMarkLength;

                        var isRotated = xScale.labelRotation !== 0;
                        var yTickStart = xScale.top;
                        var yTickEnd = xScale.top + tl;
                        var chartArea = chart.chartArea;

                        // use the saved ticks
                        var maxTicks = xScale.options.ticks._maxTicksLimit - 1;
                        var ticksPerVisibleTick = xScale.ticks.length / maxTicks;

                        // chart.js uses an integral skipRatio - this causes all the fractional ticks to be accounted for between the last 2 labels
                        // we use a fractional skipRatio
                        var ticksCovered = 0;
                        helpers.each(xScale.ticks, function (label, index) {
                            if (index < ticksCovered)
                                return;

                            ticksCovered += ticksPerVisibleTick;

                            // chart.js has already drawn these 2
                            if (index === 0 || index === (xScale.ticks.length - 1))
                                return;

                            // copy of chart.js code
                            var xLineValue = this.getPixelForTick(index);
                            var xLabelValue = this.getPixelForTick(index, this.options.gridLines.offsetGridLines);

                            if (this.options.gridLines.display) {
                                this.ctx.lineWidth = this.options.gridLines.lineWidth;
                                this.ctx.strokeStyle = this.options.gridLines.color;

                                xLineValue += helpers.aliasPixel(this.ctx.lineWidth);

                                // Draw the label area
                                this.ctx.beginPath();

                                if (this.options.gridLines.drawTicks) {
                                    this.ctx.moveTo(xLineValue, yTickStart);
                                    this.ctx.lineTo(xLineValue, yTickEnd);
                                }

                                // Draw the chart area
                                if (this.options.gridLines.drawOnChartArea) {
                                    this.ctx.moveTo(xLineValue, chartArea.top);
                                    this.ctx.lineTo(xLineValue, chartArea.bottom);
                                }

                                // Need to stroke in the loop because we are potentially changing line widths & colours
                                this.ctx.stroke();
                            }

                            if (this.options.ticks.display) {
                                this.ctx.save();
                                this.ctx.translate(xLabelValue + this.options.ticks.labelOffset, (isRotated) ? this.top + 12 : this.options.position === 'top' ? this.bottom - tl : this.top + tl);
                                this.ctx.rotate(helpers.toRadians(this.labelRotation) * -1);
                                this.ctx.font = tickLabelFont;
                                this.ctx.textAlign = (isRotated) ? 'right' : 'center';
                                this.ctx.textBaseline = (isRotated) ? 'middle' : this.options.position === 'top' ? 'bottom' : 'top';
                                this.ctx.fillText(label, 0, 0);
                                this.ctx.restore();
                            }
                        }, xScale);
                    }
                };
            }
        },
    });

    useEffect(() => {
        const fetchData = async () => {
            const options = {
                method: 'get',
                url: `/api/objects/${props.sensor.oid}`,
                time: 3000
            };

            let response;
            try {
                response = await axios(options);
            } catch (e) {
                alert('Unable to connect to backend server. Make sure the backend server is running');
                console.log(e);
                return response;
            }

            if (isEmpty(response.data.data)) {
                console.log('no data received');
                return;
            }

            if (isEmpty(response.data.labels)) {
                console.log('no labels received');
                return;
            }

            let l = [];
            let d = [];

            for (let i = 0; i <= 100; i++) {
                l[i] = i;
                d[i] = 1000;
            }

            setDataCO2({
                labels: l,
                datasets: [
                    {
                        label: [props.sensor.name],
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
                    },
                    {
                        label: 'critical',
                        fill: false,
                        lineTension: 0.1,
                        backgroundColor: 'red',
                        borderColor: 'red',
                        borderCapStyle: 'butt',
                        borderDash: [],
                        borderDashOffset: 0.0,
                        borderJoinStyle: 'miter',
                        pointBorderColor: 'red',
                        pointBackgroundColor: '#fff',
                        pointBorderWidth: 1,
                        pointHoverRadius: 5,
                        pointHoverBackgroundColor: 'red',
                        pointHoverBorderColor: 'red',
                        pointHoverBorderWidth: 2,
                        pointRadius: 1,
                        pointHitRadius: 10,
                        data: d

                    }
                ]
            });
        };

        fetchData().catch((error) => {
            console.log(error);
        });

    }, [toggle, props.sensor.name, props.sensor.oid]);

    const toggleFetch = () => setToggle(!toggle);

    return (
        <div>
            <Button style={{float: 'right'}} onClick={toggleFetch}>Refresh</Button>
            <Line data={dataCO2} options={setYAxis('ppm')}/>
        </div>
    );
}