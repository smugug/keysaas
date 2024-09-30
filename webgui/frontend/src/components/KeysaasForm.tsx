import React from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';
import {createKeySaaSInstance} from '../api'
import './KeysaasForm.css'; // Import the CSS file
import { KeySaaSResource } from '../constants';

// Validation schema using Yup
// const validationSchema = Yup.object({
//   name: Yup.string().required('Name is required'),
//   keycloakUsername: Yup.string().required('Username is required'),
//   keycloakPassword: Yup.string().required('Password is required'),
//   limitsCpu: Yup.string().email('CPU limits').required('Admin Email is required'),
//   limitsMemory: Yup.string().required('Memory limits'),
//   tls: Yup.string().required('TLS Flag is required'),
// });

const KeysaasForm: React.FC = () => {
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
      limitsCpu: '500m',
      limitsMemory: '512Mi',
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

            {/* MySQL Username */}
            <div>
              <label htmlFor="spec.keysaasUsername">Keycloak Username</label>
              <Field name="spec.keysaasUsername" type="text" />
              <ErrorMessage name="spec.keysaasUsername" component="div" className="error" />
            </div>

            {/* MySQL Password */}
            <div>
              <label htmlFor="spec.keysaasPassword">Keycloak Password</label>
              <Field name="spec.keysaasPassword" type="password" />
              <ErrorMessage name="spec.keysaasPassword" component="div" className="error" />
            </div>

            {/* Keysaas Admin Email */}
            <div>
              <label htmlFor="spec.limitsCpu">Limits CPU</label>
              <Field name="spec.limitsCpu" type="text"/>
              <ErrorMessage name="spec.limitsCpu" component="div" className="error" />
            </div>

            {/* PVC Volume Name */}
            <div>
              <label htmlFor="spec.limitsMemory">Limits Memory</label>
              <Field name="spec.limitsMemory" type="text" />
              <ErrorMessage name="spec.limitsMemory" component="div" className="error" />
            </div>

            {/* Domain Name */}
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
