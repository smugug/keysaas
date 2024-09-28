import React from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';
import {createKeySaaSInstance} from '../api'
import './KeysaasForm.css'; // Import the CSS file

// Validation schema using Yup
const validationSchema = Yup.object({
  mySQLServiceName: Yup.string().required('MySQL Service Name is required'),
  mySQLUserName: Yup.string().required('MySQL Username is required'),
  mySQLUserPassword: Yup.string().required('MySQL Password is required'),
  keysaasAdminEmail: Yup.string().email('Invalid email address').required('Admin Email is required'),
  pvcVolumeName: Yup.string().required('PVC Volume Name is required'),
  domainName: Yup.string().required('Domain Name is required'),
  tls: Yup.string().required('TLS Flag is required'),
});

interface FormValues {
  mySQLServiceName: string;
  mySQLUserName: string;
  mySQLUserPassword: string;
  keysaasAdminEmail: string;
  pvcVolumeName: string;
  domainName: string;
  tls: string;
}

const KeysaasForm: React.FC = () => {
  const initialValues: FormValues = {
    mySQLServiceName: '',
    mySQLUserName: '',
    mySQLUserPassword: '',
    keysaasAdminEmail: '',
    pvcVolumeName: '',
    domainName: '',
    tls: '',
  };

  const handleSubmit = (values: FormValues) => {
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
        validationSchema={validationSchema}
        onSubmit={handleSubmit}
      >
        {({ isSubmitting }) => (
          <Form>
            {/* MySQL Service Name */}
            <div>
              <label htmlFor="mySQLServiceName">MySQL Service Name</label>
              <Field name="mySQLServiceName" type="text" />
              <ErrorMessage name="mySQLServiceName" component="div" className="error" />
            </div>

            {/* MySQL Username */}
            <div>
              <label htmlFor="mySQLUserName">MySQL Username</label>
              <Field name="mySQLUserName" type="text" />
              <ErrorMessage name="mySQLUserName" component="div" className="error" />
            </div>

            {/* MySQL Password */}
            <div>
              <label htmlFor="mySQLUserPassword">MySQL Password</label>
              <Field name="mySQLUserPassword" type="password" />
              <ErrorMessage name="mySQLUserPassword" component="div" className="error" />
            </div>

            {/* Keysaas Admin Email */}
            <div>
              <label htmlFor="keysaasAdminEmail">Keysaas Admin Email</label>
              <Field name="keysaasAdminEmail" type="email" />
              <ErrorMessage name="keysaasAdminEmail" component="div" className="error" />
            </div>

            {/* PVC Volume Name */}
            <div>
              <label htmlFor="pvcVolumeName">PVC Volume Name</label>
              <Field name="pvcVolumeName" type="text" />
              <ErrorMessage name="pvcVolumeName" component="div" className="error" />
            </div>

            {/* Domain Name */}
            <div>
              <label htmlFor="domainName">Domain Name</label>
              <Field name="domainName" type="text" />
              <ErrorMessage name="domainName" component="div" className="error" />
            </div>

            {/* TLS Flag */}
            <div>
              <label htmlFor="tls">TLS Flag</label>
              <Field as="select" name="tls">
                <option value="">Select</option>
                <option value="true">True</option>
                <option value="false">False</option>
              </Field>
              <ErrorMessage name="tls" component="div" className="error" />
            </div>

            <button type="submit" disabled={isSubmitting}>
              Create KeySaaS
            </button>
          </Form>
        )}
      </Formik>
    </div>
  );
};

export default KeysaasForm;
