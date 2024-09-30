import React from 'react';
import KeysaasForm from './KeysaasForm';
import KeysaasList from './KeysaasList';

const Index: React.FC = () => {

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

export default Index;
