import React, { useEffect, useState } from "react";
import { Empty, Form, Input, InputNumber, Select, Tooltip } from "antd";
import axios from "axios";

// Define types for OpenAPI schema
type OpenAPISchema = {
  [key: string]: {
    type: string;
    description?: string;
    enum?: string[];
  };
};

interface DynamicFormProps {
  resource: string; // Resource type, e.g., "Pod"
}

const DynamicForm: React.FC<DynamicFormProps> = ({ resource }) => {
  const [schema, setSchema] = useState<OpenAPISchema | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    // Fetch OpenAPI JSON from Kubernetes for the given resource
    const fetchSchema = async () => {
      try {
        // 
        let response = await axios.get('https://kubernetes.default.svc/openapi/v2', { timeout: 500 });
        if (!response.data || Object.keys(response.data).length == 0)
        response = await axios.get('file:///mnt/Desktop/keysaas/openapi.json', { timeout: 1000 });
        const openApi = response;
        
        // Assuming the schema for the specific resource exists
        // const resourceSchema = openApi.
        // setSchema(resourceSchema);
      } catch (error) {
        console.error("Error fetching OpenAPI schema", error);
      }
    };

    fetchSchema();
  }, [resource]);

  // Function to generate form fields based on schema
  const generateFormFields = () => {
    if (!schema) return null;

    return Object.keys(schema).map((field) => {
      const fieldInfo = schema[field];
      const { type, description, enum: enumOptions } = fieldInfo;

      let inputField;

      // Create form fields based on the field type
      switch (type) {
        case "string":
          if (enumOptions) {
            // If enums are present, create a Select component
            inputField = (
              <Select>
                {enumOptions.map((option) => (
                  <Select.Option key={option} value={option}>
                    {option}
                  </Select.Option>
                ))}
              </Select>
            );
          } else {
            inputField = <Input />;
          }
          break;
        case "integer":
          inputField = <InputNumber />;
          break;
        case "boolean":
          inputField = <Select>
            <Select.Option value={true}>True</Select.Option>
            <Select.Option value={false}>False</Select.Option>
          </Select>;
          break;
        default:
          inputField = <Input />;
      }

      // Wrap the input field with a Tooltip for hover help
      return (
        <Form.Item
          key={field}
          label={field}
          name={field}
          tooltip={description ? <Tooltip title={description}>{description}</Tooltip> : null}
        >
          {inputField}
        </Form.Item>
      );
    });
  };

  const onSubmit = (values: any) => {
    console.log("Form values:", values);
    // Logic to create resource based on form data
  };

  return (
    <Form form={form} onFinish={onSubmit} layout="vertical">
      {generateFormFields()}
      <Form.Item>
        <button type="submit">Submit</button>
      </Form.Item>
    </Form>
  );
};

export default DynamicForm;
