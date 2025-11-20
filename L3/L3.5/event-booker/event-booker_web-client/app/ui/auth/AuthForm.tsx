'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Box, Typography, Alert } from '@mui/material';
import Card from '../Card';
import Button from '../Button';
import Input from '../Input';
import { LoginRequest, RegisterRequest } from '../../lib/types';

interface AuthFormProps {
  mode: 'login' | 'register';
}

export default function AuthForm({ mode }: AuthFormProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const credentials: LoginRequest | RegisterRequest = {
        email,
        password
      };

      if (mode === 'register') {
        (credentials as RegisterRequest).name = name;
      }

      // Здесь будет вызов API для аутентификации
      // const response = await (mode === 'login' ? login(credentials) : register(credentials));
      // localStorage.setItem('token', response.token);
      // router.push('/events');

      // Заглушка для демонстрации
      console.log(mode === 'login' ? 'Login attempt' : 'Register attempt', credentials);
      router.push('/events');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card>
      <Typography variant="h5" component="h1" gutterBottom>
        {mode === 'login' ? 'Вход' : 'Регистрация'}
      </Typography>
      
      <Box component="form" onSubmit={handleSubmit} sx={{ width: '100%' }}>
        {mode === 'register' && (
          <Input
            label="Имя"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        )}
        
        <Input
          label="Email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        
        <Input
          label="Пароль"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        
        {error && (
          <Alert severity="error" sx={{ mt: 2, width: '100%' }}>
            {error}
          </Alert>
        )}
        
        <Button 
          type="submit" 
          disabled={loading}
          sx={{ mt: 2 }}
        >
          {loading ? 'Загрузка...' : (mode === 'login' ? 'Войти' : 'Зарегистрироваться')}
        </Button>
      </Box>
      
      <Box sx={{ mt: 2, textAlign: 'center' }}>
        <Typography variant="body2" color="text.secondary">
          {mode === 'login' 
            ? "Нет аккаунта? " 
            : "Уже есть аккаунт? "}
          <a 
            href={mode === 'login' ? '/auth?mode=register' : '/auth?mode=login'}
          >
            {mode === 'login' ? 'Зарегистрироваться' : 'Войти'}
          </a>
        </Typography>
      </Box>
    </Card>
  );
}