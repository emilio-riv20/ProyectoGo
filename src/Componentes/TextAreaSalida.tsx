import React from 'react';

interface TextareaSalidaProps {
  value: string;
}

const TextareaSalida: React.FC<TextareaSalidaProps> = ({ value }) => {
  return (
    <textarea
      value={value || ''}
      readOnly
      placeholder="Los resultados apareceran acá"
    />
  );
};

export default TextareaSalida;
