import React from 'react';

interface TextareaEntradaProps {
  value: string;
  onChange: (value: string) => void;
}

const TextareaEntrada: React.FC<TextareaEntradaProps> = ({ value, onChange }) => {
  return (
    <textarea
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder="Escribe comandos acá"
    />
  );
};

export default TextareaEntrada;
