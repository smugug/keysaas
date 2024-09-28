import React from 'react';
import KeysaasForm from './components/KeysaasForm';
import KeysaasList from './components/KeysaasList';
import './App.css';

const App: React.FC = () => {

  localStorage.setItem('token', 'eyJhbGciOiJSUzI1NiIsImtpZCI6ImdTaTl5aVhHT0tGTFlETXI1UEtYZVBUdmlOYVprN0ZodWRadlUzWGI3ckkifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzMxMDk0NDY0LCJpYXQiOjE3Mjc0OTQ0NjQsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiNTBiNjBkNDQtYmY4OC00MTVjLWE0MGQtNjcyZmM3ODViYjBhIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJjdXN0b21lcjIiLCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoid2ViZ3VpLWFjY291bnQiLCJ1aWQiOiJiZmFhOWI0YS1jZWRjLTQ0ZDUtOGIwMi00NTA2YTk0ZTI4ZTYifX0sIm5iZiI6MTcyNzQ5NDQ2NCwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmN1c3RvbWVyMjp3ZWJndWktYWNjb3VudCJ9.C_iTdbQm-xHrHdrIdJodizz_YO1iFFVaZzgOFfCsxSBNV6MFyE820FlpFona9hDYc2zMUujbqri55kJTfqH8NTktlvdkKdJ2DR8--vh7NSeUJORBhMahMhxG8XNiltr_HwXo8ad3UWTcWAHRgH90b2jh5mjXLAA-YDIOiHfZ5I6V6Ib9aWWsQr9OooOuYq4AV94tmSvhs6jMvhqnfk28NpfUq9YIDsG7h5RxaKoHVP8EuHZbVFkAI5kdvGM9ZRAGQnqgOfm1p7a9yKRtjcdyR2GDTd_xqgnocfCq0_wzu9LPtAQ_fbk_X7P10SM9yUB2vvEl2wZo5QsqvUCD8CUibg');
  return (
    <div className="container">
      <div className="form">
        <KeysaasForm/>
      </div>
      <div className="list">
        <KeysaasList/>
      </div>
    </div>

  );
};

export default App;
