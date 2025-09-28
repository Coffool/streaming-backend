import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import ProtectedRoute from "@/components/ProtectedRoute";
import Login from "@/pages/Login";
import UserDashboard from "@/pages/UserDashboard";
import ArtistDashboard from "@/pages/ArtistDashboard";


export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      
      <Route path="/dashboard/home" element={
        <ProtectedRoute>
          <UserDashboard />
        </ProtectedRoute>
      } />
      
      <Route path="/dashboard/artist" element={
        <ProtectedRoute requiredRole="artist">
          <ArtistDashboard />
        </ProtectedRoute>
      } />
      
      <Route path="/unauthorized" element={<div>No tienes permisos para acceder a esta página</div>} />
    </Routes>
  );
}
