// src/pages/Login.jsx
import React, { useRef, useState, useLayoutEffect } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";

import Button from "@/components/ui/button";
import Input from "@/components/ui/input";
import Label from "@/components/ui/label";
import Card, { CardHeader, CardContent, CardFooter, CardTitle, CardDescription } from "@/components/ui/card";
import Checkbox from "@/components/ui/checkbox";
import Spinner from "@/components/ui/spinner";
import { Music } from "lucide-react";

import { useAuth } from "@/contexts/AuthContexts";

/**
 * Login page con:
 * - tabs con indicador animado (medido en px)
 * - animación de entrada para logo + card usando framer-motion
 */

export default function Login() {
  const navigate = useNavigate();
  const { login, register} = useAuth();
  const [active, setActive] = useState("login"); // "login" | "register"
  const [authError, setAuthError] = useState("");

  // refs y estado para el indicador
  const tabsRef = useRef(null);
  const [indicatorWidth, setIndicatorWidth] = useState(0); // px
  const [indicatorX, setIndicatorX] = useState(0); // px

  // buffer extra para cubrir redondeos 
  const BUFFER_PX = 6;

  // medir y actualizar indicador (usa useLayoutEffect para leer layout antes de paint)
  useLayoutEffect(() => {
    if (!tabsRef.current) return;
    const el = tabsRef.current;

    const resize = () => {
      const w = el.clientWidth;
      // asumimos dos tabs iguales -> indicador ocupa la mitad (+ buffer)
      const half = Math.round(w / 2) + BUFFER_PX;
      setIndicatorWidth(half);

      // posición según active (posición en px)
      setIndicatorX(active === "login" ? 0 : Math.round(w / 2));
    };

    // initial
    resize();

    // ResizeObserver para cambios de layout internos
    const ro = new ResizeObserver(resize);
    ro.observe(el);

    // window events (resize/zoom/orientation)
    window.addEventListener("resize", resize);
    window.addEventListener("orientationchange", resize);

    return () => {
      ro.disconnect();
      window.removeEventListener("resize", resize);
      window.removeEventListener("orientationchange", resize);
    };
  }, [tabsRef, active]); // recalcula cuando active cambie o el ref esté listo

  // Login form
  const {
    register: registerLogin,
    handleSubmit: handleSubmitLogin,
    formState: { errors: loginErrors, isSubmitting: loginSubmitting },
  } = useForm();

  // Register form
  const {
    register: registerReg,
    handleSubmit: handleSubmitReg,
    watch: watchReg,
    formState: { errors: regErrors, isSubmitting: regSubmitting },
  } = useForm();

  const onLogin = async (values) => {
    try {
      setAuthError("");
      // Preparar datos para el backend - el endpoint espera "identifier" (username o email)
      const credentials = {
        identifier: values.email,
        password: values.password,
      };

      const response = await login(credentials);
      
      // Redirigir según el rol del usuario
      if (response.user.role === 'artist') {
        navigate("/dashboard/artist");
      } else {
        navigate("/dashboard/home");
      }
    } catch (err) {
      console.error("login error", err);
      setAuthError(err.message || "Error durante el login");
    }
  };

  const onRegister = async (values) => {
    if (values.password !== values.confirmPassword) {
      setAuthError("Las contraseñas no coinciden");
      return;
    }
    
    try {
      setAuthError("");
      const userData = {
        username: values.username,
        email: values.email,
        password: values.password,
        birthdate: values.birthdate,
        is_artist: values.isArtist || false,
      };

      await register(userData);
      alert("Registro exitoso. Por favor inicia sesión.");
      setActive("login");
    } catch (err) {
      console.error("register error", err);
      setAuthError(err.message || "Error durante el registro");
    }
  };
  // manejo teclado para cambiar tabs (Left/Right)
  const onKeyDownTabs = (e) => {
    if (e.key === "ArrowLeft") setActive("login");
    if (e.key === "ArrowRight") setActive("register");
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-white via-slate-50 to-gray-100 flex items-center justify-center px-6 py-12">
      <div className="w-full max-w-lg">
        {/* LOGO + tagline: animación de entrada */}
        <motion.div
          initial={{ opacity: 0, y: -8, scale: 0.99 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ duration: 0.45, ease: "easeOut", delay: 0.06 }}
          className="text-center mb-6"
        >
          <div className="inline-flex items-center gap-3 mb-2 justify-center">
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ duration: 0.45, delay: 0.12 }}
              className="w-12 h-12 bg-gradient-to-br from-indigo-600 to-pink-500 rounded-xl flex items-center justify-center"
            >
              <Music className="w-6 h-6 text-white" />
            </motion.div>

            <motion.span
              initial={{ opacity: 0, y: -4 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.45, delay: 0.14 }}
              className="font-bold text-2xl bg-gradient-to-r from-indigo-600 to-pink-500 bg-clip-text text-transparent"
            >
              VibeStream
            </motion.span>
          </div>
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.45, delay: 0.18 }}
            className="text-sm text-gray-500"
          >
            Discover, create and share your musical journey
          </motion.p>
        </motion.div>

        {/* Card con animación de entrada (slide-up + fade + scale suave) */}
        <motion.div
          initial={{ opacity: 0, y: 20, scale: 0.995 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ duration: 0.6, ease: "easeOut", delay: 0.22 }}
        >
          <Card className="shadow-lg">
            {/* Header con tabs control */}
            <CardHeader className="px-4 pt-4 pb-0">
              <div
                role="tablist"
                aria-label="Auth tabs"
                ref={tabsRef}
                onKeyDown={onKeyDownTabs}
                tabIndex={0}
                className="relative w-full max-w-md mx-auto bg-gray-100 rounded-md p-1"
              >
                {/* sliding background animado en px (left: 0; animamos translateX) */}
                <motion.div
                  animate={{ x: indicatorX }}
                  transition={{ type: "spring", stiffness: 300, damping: 30 }}
                  className="pointer-events-none absolute top-0 h-full rounded-md bg-gradient-to-r from-indigo-500 to-purple-600 shadow-lg"
                  style={{
                    width: `${indicatorWidth}px`,
                    left: 0,
                    willChange: "transform",
                  }}
                  aria-hidden
                />

                <div className="relative z-10 flex">
                  <button
                    role="tab"
                    aria-selected={active === "login"}
                    onClick={() => setActive("login")}
                    className={`flex-1 z-20 px-4 py-2 text-sm font-medium rounded-md focus:outline-none transition-colors ${
                      active === "login" ? "text-white" : "text-gray-600"
                    }`}
                  >
                    Sign In
                  </button>

                  <button
                    role="tab"
                    aria-selected={active === "register"}
                    onClick={() => setActive("register")}
                    className={`flex-1 z-20 px-4 py-2 text-sm font-medium rounded-md focus:outline-none transition-colors ${
                      active === "register" ? "text-white" : "text-gray-600"
                    }`}
                  >
                    Sign Up
                  </button>
                </div>
              </div>
            </CardHeader>

            <CardContent className="pt-6">
              {active === "login" ? (
                <form onSubmit={handleSubmitLogin(onLogin)}>
                  <div className="space-y-4">
                    <CardTitle className="text-center text-2xl">Welcome back!</CardTitle>
                    <CardDescription className="text-center">Sign in to continue</CardDescription>

                    <div>
                      <Label htmlFor="login-email">Email</Label>
                      <Input
                        id="login-email"
                        type="email"
                        placeholder="you@example.com"
                        {...registerLogin("email", { required: "Email requerido" })}
                        className="mt-1"
                      />
                      {loginErrors.email && <p className="text-xs text-red-500 mt-1">{loginErrors.email.message}</p>}
                    </div>

                    <div>
                      <Label htmlFor="login-password">Password</Label>
                      <Input
                        id="login-password"
                        type="password"
                        {...registerLogin("password", { required: "Password requerido" })}
                        className="mt-1"
                      />
                      {loginErrors.password && <p className="text-xs text-red-500 mt-1">{loginErrors.password.message}</p>}
                    </div>
                  </div>

                  <CardFooter className="flex flex-col space-y-3 pt-4">
                    <Button type="submit" variant="gradient" className="w-full" loading={loginSubmitting}>
                      {loginSubmitting ? <Spinner /> : "Sign In"}
                    </Button>
                  </CardFooter>
                </form>
              ) : (
                <form onSubmit={handleSubmitReg(onRegister)}>
                  <div className="space-y-4">
                    <CardTitle className="text-center text-2xl">Create account</CardTitle>
                    <CardDescription className="text-center">Join VibeStream — it’s free</CardDescription>  

                    <div>
                      <Label htmlFor="reg-username">Username</Label>
                      <Input id="reg-name" {...registerReg("username", { required: "Username required" })} className="mt-1" />
                      {regErrors.name && <p className="text-xs text-red-500 mt-1">{regErrors.name.message}</p>}
                    </div>

                    <div>
                      <Label htmlFor="reg-email">Email</Label>
                      <Input id="reg-email" type="email" {...registerReg("email", { required: "Email requerido" })} className="mt-1" />
                      {regErrors.email && <p className="text-xs text-red-500 mt-1">{regErrors.email.message}</p>}
                    </div>

                    <div>
                      <Label htmlFor='reg-birthdate'>Birthdate</Label>
                      <Input id='reg-birthdate' type='date' {...registerReg("birthdate", { required: "Birthdate required",  })} className="mt-1" />
                      {regErrors.birthdate && <p className="text-xs text-red-500 mt-1">{regErrors.birthdate.message}</p>}
                    
                    </div>

                    <div>
                      <Label htmlFor="reg-password">Password</Label>
                      <Input
                        id="reg-password"
                        type="password"
                        {...registerReg("password", { required: "Password requerido", minLength: { value: 8, message: "Min 8 chars" } })}
                        className="mt-1"
                      />
                      {regErrors.password && <p className="text-xs text-red-500 mt-1">{regErrors.password.message}</p>}
                    </div>

                    <div>
                      <Label htmlFor="reg-confirm">Confirm password</Label>
                      <Input id="reg-confirm" type="password" {...registerReg("confirmPassword", { required: "Confirma contraseña" })} className="mt-1" />
                      {regErrors.confirmPassword && <p className="text-xs text-red-500 mt-1">{regErrors.confirmPassword.message}</p>}
                    </div>

                    <div className="flex items-center gap-2">
                      <Checkbox id="artist-register" {...registerReg("isArtist")} />
                      <Label htmlFor="artist-register" className="text-sm">I'm an artist</Label>
                    </div>
                  </div>

                  <CardFooter className="flex flex-col space-y-3 pt-4">
                    <Button type="submit" variant="gradient" className="w-full" loading={regSubmitting}>
                      {regSubmitting ? <Spinner /> : "Create Account"}
                    </Button>
                  </CardFooter>
                </form>
              )}
            </CardContent>
          </Card>
        </motion.div>
      </div>
    </div>
  );
}
