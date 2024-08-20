import React from 'react';

interface BotonesProps {
  onArchivoCargado: (contenido: string) => void;
  onEjecutar: () => void;
}

const Botones: React.FC<BotonesProps> = ({ onArchivoCargado, onEjecutar }) => {
  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        const contenido = e.target?.result as string;
        onArchivoCargado(contenido);
      };
      reader.readAsText(file);
    }
  };

  return (
    <div>
      <input type="file" onChange={handleFileChange} accept=".mia" />
      <button onClick={onEjecutar}>Ejecutar</button>
    </div>
  );
};

export default Botones;
