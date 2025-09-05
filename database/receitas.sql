DROP TABLE IF EXISTS receitas;
CREATE TABLE receitas (
  id         INT AUTO_INCREMENT NOT NULL,
  receita      VARCHAR(255) NOT NULL,
  tipo         ENUM('Peixe', 'Carne', 'Veg') NOT NULL,
  proteina     ENUM('Frango', 'Peru', 'Porco', 'Vaca', 'Borrego', 'Soja', 'Grão de Bico', 'Peixe Variado', 'Bacalhau', 'Frutos do Mar/Moluscos', 'Atum', 'Salmão') NOT NULL,
  domingo      ENUM('true', 'false') NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO receitas
  (receita, tipo, proteina, domingo)
VALUES
  ('Feijoada', 'Carne', 'Porco', 'true'),
  ('Bacalhau à Braz', 'Peixe', 'Bacalhau', 'false'),
  ('Carne estufada', 'Carne', 'Vaca', 'false'),
  ('Bifes de Frango', 'Carne', 'Frango', 'false'),
  ('Arroz de Polvo', 'Peixe', 'Frutos do Mar/Moluscos', 'true'),
  ('Lasanha de Atum', 'Peixe', 'Atum', 'false'),
  ('Carne Tropical', 'Carne', 'Porco', 'true'),
  ('Filetes de Peixe', 'Peixe', 'Peixe Variado', 'false'),
  ('Massa com Salmão', 'Peixe', 'Salmão', 'false'),
  ('Bacalhau com Natas', 'Peixe', 'Bacalhau', 'true'),
  ('Carne Guisada', 'Carne', 'Porco', 'false'),
  ('Robalo Grelhado', 'Peixe', 'Peixe Variado', 'false'),
  ('Caril de Lulas', 'Peixe', 'Frutos do Mar/Moluscos', 'true'),
  ('Bifes Chouriço', 'Carne', 'Frango', 'false'),
  ('Arroz de Nozes', 'Carne', 'Frango', 'false'),
  ('Jardineira', 'Carne', 'Frango','false'),
  ('Lombinhos na Frigideira', 'Carne', 'Porco', 'false'),
  ('Peixe à Lumbo', 'Peixe', 'Peixe Variado', 'true')
