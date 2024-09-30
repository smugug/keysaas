import React from 'react';
import { Routes, Route, BrowserRouter } from 'react-router-dom';
import './App.css';
import KeysaasDetail from './components/KeysaasDetail';
import Index from './components/Index';

const App: React.FC = () => {

  localStorage.setItem('token', 'eyJhbGciOiJSUzI1NiIsImtpZCI6InFhRTR2N1VvZFg2c1Fha2FISHo1UVRpS1VmRzI4elNzMXp6elpNU2Zhb2cifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzMxMTc3Nzk0LCJpYXQiOjE3Mjc1Nzc3OTQsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiN2E2OWZiOWMtNGE5Ny00NjM0LTkzNWItMjM3NGI1NzczZWRiIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJjdXN0b21lcjIiLCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoid2ViZ3VpLWFjY291bnQiLCJ1aWQiOiIwNTMxZTU0OS05ZTllLTQwMTYtOGMyYS1lOTQ4ZGJmODk3NzEifX0sIm5iZiI6MTcyNzU3Nzc5NCwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmN1c3RvbWVyMjp3ZWJndWktYWNjb3VudCJ9.lxwvVek26xcdVv2u4QnYliBtccxum7uMMJEF3ug2yf75Fl7VPI4p-1w4U3M7HdBIqyItlyTPCAd3gcfK_gl4lRFxw-V-CPa-PD7xfctUP0OaxoVkn36JOk2NZCZ7AM86nvkII_pOI3QWYOAadz-A2QrK-Wj1pV-XsML2LQaM9oDBvq9QIyeFuOwYqnKte5j5Y0_WFq0M2GYY0fUKLvJg2Qq9-RYMcQBi_Jisl-7mQ-MFPqFyOL_NFFdHBFadxqV2Aqghc2_wkFOGFtpep47az2l6umYXoUN37OqISlAHm7xlvnO9HBvOOWTjruq7QMO1jzcJhvnHUPftuaKXJNwKcA');
  return (
    <BrowserRouter>
      <Routes>
        <Route 
          path="/" 
          element={<Index/>} 
        />
        <Route
          path="/keysaas/:name"
          element={<KeysaasDetail/>}
        />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
