import React, {useEffect, useState} from 'react';
import { Line } from "react-chartjs-2";
import { Chart as ChartJS, LineElement, CategoryScale, LinearScale, PointElement, Title, Tooltip, Legend } from 'chart.js';

ChartJS.register(LineElement, CategoryScale, LinearScale, PointElement, Title, Tooltip, Legend);

const LineChart = () => {
    const nameMapping = {
        memory: 0,
        cpu: 1
    }
    const [chartData, setChartData] = useState({
        labels: [],
        datasets: [
            {
                label: "Memory Usage %",
                data: [],
                borderColor: "blue",
                backgroundColor: "rgba(0, 0, 255, 0.2)",
                borderWidth: 2,
                pointRadius: 5,
            },
            {
                label: "CPU Usage %",
                data: [],
                borderColor: "red",
                backgroundColor: "rgba(0, 0, 255, 0.2)",
                borderWidth: 2,
                pointRadius: 5,
            },
        ],
    });

    useEffect(() => {
        const eventSource = new EventSource('http://localhost:8001/events');
        eventSource.onmessage = (event) => {
            const newData = JSON.parse(event.data);
            console.log(newData);

            setChartData((prevData) => {
                const updatedLabels = [...prevData.labels, newData.createdAt].slice(-5);

                const datasets = []
                for (let i = 0; i < newData.value.length; i++) {
                    let obj = newData.value[i]
                    let indexValue = nameMapping[obj.labelKey]

                    datasets.push(
                        {
                            ...prevData.datasets[indexValue], 
                            data: [...prevData.datasets[indexValue].data.slice(-4), obj.labelValue]
                        }
                    )
                }
                
                console.log(datasets)
                return {
                    ...prevData,
                    labels: updatedLabels,
                    datasets: datasets
                };
            });
        };

        return () => {
            eventSource.close()
        };
    }, []);

    const options = {
        responsive: true, 
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: true,
                position: 'top'
            },
        },
        scales: {
            x: {
                grid: {
                    display: false
                },
            },
            y: {
                beginAtZero: true,
                min: 0,
                max: 100
            }
        }
    };

    return (
        <div style={{ width: '600px', height: '400px' }}>
            <Line data={chartData} options={options} />
        </div>
    );
};

export default LineChart;
