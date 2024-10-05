import React, { useState } from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
// import * as Yup from 'yup';
import {createKeySaaSInstance} from '../api'
import './KeysaasForm.css'; // Import the CSS file
import { KeySaaSResource } from '../constants';

// Validation schema using Yup
// const validationSchema = Yup.object({
//   name: Yup.string().required('Name is required'),
//   keycloakUsername: Yup.string().required('Username is required'),
//   keycloakPassword: Yup.string().required('Password is required'),
//   tls: Yup.string().required('TLS Flag is required'),
// });

const KeysaasForm: React.FC = () => {
  const [useExternalDb, setUseExternalDb] = useState(false);
  const initialValues: KeySaaSResource = {
    apiVersion: 'keysaascontroller.keysaas/v1',
    kind: 'Keysaas',
    metadata: {
      name: '',
      namespace: 'customer2',
    },
    spec: {
      domainName: '',
      keysaasPassword: '',
      keysaasUsername: '',
      requestsCpu: '500m',
      requestsMemory: '512Mi',
      limitsCpu: '500m',
      limitsMemory: '512Mi',
      scalingThreshold: '90',
      minInstances: '1',
      maxInstances: '2',
      postgresUri: '',
      tls: 'true',
    },
  }

  const handleSubmit = (values: KeySaaSResource) => {
    console.log('Submitted values:', values);
    try {
      const data = createKeySaaSInstance(values);
      console.log(data)
    } catch (error) {
      console.error('Error creating Keysaas instances:', error);
    }
  };

  // const validationSchema = Yup.object().shape({
  //   keysaasUsername: Yup.string().required("Required"),
  //   keysaasPassword: Yup.string().required("Required"),
  //   requestsMemory: Yup.string(),
  //   requestsCpu: Yup.string(),
  //   limitsMemory: Yup.string(),
  //   limitsCpu: Yup.string(),
  //   scalingThreshold: Yup.string(),
  //   minInstances: Yup.string(),
  //   maxInstances: Yup.string(),
  //   postgresUri: useExternalDb
  //     ? Yup.string().required("Required if using an external database")
  //     : Yup.string(),
  //   domainName: Yup.string(),
  //   tls: Yup.string().required("Required"),
  // });


  return (  
    <div>
      <h2>Create KeySaaS Instance</h2>
      <Formik
        initialValues={initialValues}
        // validationSchema={validationSchema}
        onSubmit={handleSubmit}
      >
        {({ isSubmitting }) => (
          <Form>
            {/* MySQL Service Name */}
            <div>
              <label htmlFor="metadata.name">Name</label>
              <Field name="metadata.name" type="text" />
              <ErrorMessage name="metadata.name" component="div" className="error" />
            </div>
            <div className="inline-field">
            <div>
              <label htmlFor="spec.keysaasUsername">Keycloak Username</label>
              <Field name="spec.keysaasUsername" type="text" />
              <ErrorMessage name="spec.keysaasUsername" component="div" className="error" />
            </div>

            <div>
              <label htmlFor="spec.keysaasPassword">Keycloak Password</label>
              <Field name="spec.keysaasPassword" type="password" />
              <ErrorMessage name="spec.keysaasPassword" component="div" className="error" />
            </div>
            </div>
            <div className="inline-field">
            <div>
              <label htmlFor="spec.requestsMemory">Requests Memory</label>
              <Field name="spec.requestsMemory" type="text" />
              <ErrorMessage name="spec.requestsMemory" component="div" className="error" />
            </div>

            <div>
              <label htmlFor="spec.requestsCpu">Requests CPU</label>
              <Field name="spec.requestsCpu" type="text" />
              <ErrorMessage name="spec.requestsCpu" component="div" className="error" />
            </div>
            </div>
            <div className="inline-field">
            <div>
              <label htmlFor="spec.limitsMemory">Limits Memory</label>
              <Field name="spec.limitsMemory" type="text" />
              <ErrorMessage name="spec.limitsMemory" component="div" className="error" />
            </div>

            <div>
              <label htmlFor="spec.limitsCpu">Limits CPU</label>
              <Field name="spec.limitsCpu" type="text" />
              <ErrorMessage name="spec.limitsCpu" component="div" className="error" />
            </div>
            </div>
            <div>
              <label htmlFor="spec.scalingThreshold">Scaling Threshold</label>
              <Field name="spec.scalingThreshold" type="text" />
              <ErrorMessage name="spec.scalingThreshold" component="div" className="error" />
            </div>
            <div className="inline-field">
            <div>
              <label htmlFor="spec.minInstances">Minimum Replicas</label>
              <Field name="spec.minInstances" type="text" />
              <ErrorMessage name="spec.minInstances" component="div" className="error" />
            </div>

            <div>
              <label htmlFor="spec.maxInstances">Maximum Replicas</label>
              <Field name="spec.maxInstances" type="text" />
              <ErrorMessage name="spec.maxInstances" component="div" className="error" />
            </div>
            </div>
            <div className="inline-field">
            <div>
              <label htmlFor="useExternalDb">Use External Database</label>
              <Field as="select" name="useExternalDb" onChange={(e: React.ChangeEvent<HTMLSelectElement>) => setUseExternalDb(e.target.value === "true")}>
                <option value="false">False</option>
                <option value="true">True</option>
              </Field>
              <ErrorMessage name="tls" component="div" className="error" />
            </div>

            <div>
              <label htmlFor="postgresUri">Postgres URI</label>
              <Field name="postgresUri" type="text" disabled={!useExternalDb} />
              <ErrorMessage name="postgresUri" component="div" className="error" />
            </div>

            </div>
            <div className="inline-field">
            <div>
              <label htmlFor="spec.domainName">Domain Name</label>
              <Field name="spec.domainName" type="text"/>
              <ErrorMessage name="spec.domainName" component="div" className="error" />
            </div>

            {/* TLS Flag */}
            <div>
              <label htmlFor="spec.tls">TLS Flag</label>
              <Field as="select" name="spec.tls">
                <option value="true">True</option>
                <option value="">False</option>
              </Field>
              <ErrorMessage name="tls" component="div" className="error" />
            </div>
            </div>  
            <button type="submit" disabled={isSubmitting} onClick={()=>{

              }}>
              Create KeySaaS
            </button>
          </Form>
        )}
      </Formik>
    </div>
  );
};

export default KeysaasForm;
