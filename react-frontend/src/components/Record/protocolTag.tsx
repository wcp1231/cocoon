import styled from "@emotion/styled";

const ProtocolTag = styled.span`
  color: #fff;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.25rem 0.4rem;
  border-radius: 4px;
  
  background: ${props => props.color};;
`;

export default function(protocol: string) {
  let color = '#00000033';
  if (protocol === 'Dubbo') {
    color = '#2BA3DE';
  }
  if (protocol === 'HTTP') {
    color = '#0D6EFD';
  }
  if (protocol === 'Redis') {
    color = '#a51f17';
  }
  if (protocol === 'Mongo') {
    color = '#116149';
  }
  if (protocol === 'Mysql') {
    color = '#3E6E93';
  }
  return (<ProtocolTag color={color}>{protocol}</ProtocolTag>)
}