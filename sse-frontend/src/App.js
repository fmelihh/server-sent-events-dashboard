import React from "react";
import logo from './logo.svg';
import './App.css';
import LineChart from './components/LineChart';

function App() {
  return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        flexDirection: 'column'
      }}>
          <h2>Computer Health Dashboard</h2>
          <LineChart />
      </div>
  );
}

export default App;
