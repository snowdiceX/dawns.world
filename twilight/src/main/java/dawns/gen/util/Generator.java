package dawns.gen.util;

import java.io.File;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang.ObjectUtils;
import org.apache.velocity.VelocityContext;
import org.mybatis.generator.api.MyBatisGenerator;
import org.mybatis.generator.config.Configuration;
import org.mybatis.generator.config.xml.ConfigurationParser;
import org.mybatis.generator.internal.DefaultShellCallback;

import dawns.gen.util.mybatis.JdbcUtil;

/**
 * 代码生成类
 */
@SuppressWarnings("all")
public class Generator {

	// generatorConfig模板路径
	private static String generatorConfig_vm = "/template/generatorConfig.vm";
	// Service模板路径
	private static String service_vm = "/template/Service.vm";
	// ServiceMock模板路径
	private static String serviceMock_vm = "/template/ServiceMock.vm";
	// ServiceImpl模板路径
	private static String serviceImpl_vm = "/template/ServiceImpl.vm";

	// ControllerImpl模板路径
	private static String controller_vm = "/template/Controller.vm";

	/**
	 * 根据模板生成generatorConfig.xml文件
	 * @param jdbcDriver   驱动路径
	 * @param jdbcUrl      链接
	 * @param jdbcUsername 帐号
	 * @param jdbcPassword 密码
	 * @param module        项目模块
	 * @param database      数据库
	 * @param tablePrefix  表前缀
	 * @param packageName  包名
	 */
	public static void generator(
			String jdbcDriver,
			String jdbcUrl,
			String jdbcUsername,
			String jdbcPassword,
			String module,
			String database,
			String tablePrefix,
			String packageName,
			String autoKey) throws Exception{

		String os = System.getProperty("os.name");
//		String targetProject = module + "/" + module + "-dao";
		String targetProject = module;
		String basePath = Generator.class.getResource("/").getPath().replace("/target/classes/", "").replace(targetProject, "");
		if (os.toLowerCase().startsWith("win")) {
			generatorConfig_vm = Generator.class.getResource(generatorConfig_vm).getPath().replaceFirst("/", "");
			service_vm = Generator.class.getResource(service_vm).getPath().replaceFirst("/", "");
			serviceMock_vm = Generator.class.getResource(serviceMock_vm).getPath().replaceFirst("/", "");
			serviceImpl_vm = Generator.class.getResource(serviceImpl_vm).getPath().replaceFirst("/", "");
			basePath = basePath.replaceFirst("/", "");
		} else {
			generatorConfig_vm = Generator.class.getResource(generatorConfig_vm).getPath();
			service_vm = Generator.class.getResource(service_vm).getPath();
			serviceMock_vm = Generator.class.getResource(serviceMock_vm).getPath();
			serviceImpl_vm = Generator.class.getResource(serviceImpl_vm).getPath();
		}

		String generatorConfigXml = Generator.class.getResource("/").getPath().replace("/target/classes/", "") + "/src/main/resources/generatorConfig.xml";
		targetProject = basePath + targetProject;
		File mapperXmlPath = new File(targetProject + "/src/main/resources/gen/");
		System.out.println("mapperXmlPath: "+mapperXmlPath.getAbsolutePath());
		if(!mapperXmlPath.exists()) {
			mapperXmlPath.mkdirs();
		}
		
		String tableSql = "SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '" + database + "' AND table_name LIKE '" + tablePrefix + "_%';";

//		String tableSql = "SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '" + database + "' AND table_name ='article_info';";

		System.out.println("========== 开始生成generatorConfig.xml文件 ==========");
		List<Map<String, Object>> tables = new ArrayList<>();
		try {
			VelocityContext context = new VelocityContext();
			Map<String, Object> table;

			// 查询定制前缀项目的所有表
			JdbcUtil jdbcUtil = new JdbcUtil(jdbcDriver, jdbcUrl, jdbcUsername, jdbcPassword);
			List<Map> result = jdbcUtil.selectByParams(tableSql, null);
			for (Map map : result) {
				System.out.println(map.get("TABLE_NAME"));
				table = new HashMap<>(3);
				table.put("table_name", map.get("TABLE_NAME"));
				table.put("model_name", StringUtil.lineToHump(ObjectUtils.toString(map.get("TABLE_NAME"))));
				String columnSql = "SELECT table_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = '" + database + "' AND table_name = '" + map.get("TABLE_NAME") +"' AND column_name='"+autoKey+ "';";
				List<Map> columnResult =jdbcUtil.selectByParams(columnSql,null);
				if(columnResult.size()>0){
					table.put("has_autoKey", "true");
				}else{
					table.put("has_autoKey", "false");
				}
				tables.add(table);
			}
			jdbcUtil.release();

			context.put("tables", tables);
			context.put("autoKey", autoKey);
			context.put("generator_javaModelGenerator_targetPackage", packageName + ".gen.dao.model");
			context.put("generator_sqlMapGenerator_targetPackage",  "mapper.dao");
			context.put("generator_javaClientGenerator_targetPackage", packageName + ".gen.dao.mapper");
			context.put("targetProject", targetProject);
			context.put("targetProject_sqlMap", targetProject);
			context.put("generator_jdbc_password", jdbcPassword);
			VelocityUtil.generate(generatorConfig_vm, generatorConfigXml, context);
			// 删除旧代码
			deleteDir(new File(targetProject + "/src/main/java/" + packageName.replaceAll("\\.", "/") + "/gen/dao/model"));
			deleteDir(new File(targetProject + "/src/main/java/" + packageName.replaceAll("\\.", "/") + "/gen/dao/mapper"));
			deleteDir(new File(targetProject + "/src/main/resources/gen/mapper/dao"));
		} catch (Exception e) {
			e.printStackTrace();
		}
		System.out.println("========== 结束生成generatorConfig.xml文件 ==========");

		System.out.println("========== 开始运行MybatisGenerator ==========");
		List<String> warnings = new ArrayList<>();
		File configFile = new File(generatorConfigXml);
		ConfigurationParser cp = new ConfigurationParser(warnings);
		Configuration config = cp.parseConfiguration(configFile);
		DefaultShellCallback callback = new DefaultShellCallback(true);
		MyBatisGenerator myBatisGenerator = new MyBatisGenerator(config, callback, warnings);
		myBatisGenerator.generate(null);
		for (String warning : warnings) {
			System.out.println(warning);
		}
		System.out.println("========== 结束运行MybatisGenerator ==========");
		System.out.println("========== 删除generatorConfig.xml文件 ==========");
		configFile.delete();
		System.out.println("========== 开始生成Service ==========");
		String ctime = new SimpleDateFormat("yyyy/M/d").format(new Date());
		String servicePath = targetProject + "/src/main/java/" + packageName.replaceAll("\\.", "/") + "/gen/service";
		String serviceImplPath = servicePath+ "/impl";
		String serviceMockPath = servicePath+ "/mock";
		String controllerPath = targetProject + "/src/main/java/" + packageName.replaceAll("\\.", "/") + "/gen/controller";
		for (int i = 0; i < tables.size(); i++) {
			String model = StringUtil.lineToHump(ObjectUtils.toString(tables.get(i).get("table_name")));
			String service = servicePath + "/" + model + "Service.java";
			String serviceImpl = serviceImplPath + "/" + model + "ServiceImpl.java";
			String serviceMock = serviceMockPath + "/" + model + "ServiceMock.java";
			String controller = controllerPath + "/" + model + "RestController.java";
			// 生成service
			File serviceFile = new File(service);
			if (!serviceFile.exists()) {
				VelocityContext context = new VelocityContext();
				context.put("package_name", packageName);
				context.put("model", model);
				context.put("ctime", ctime);
				VelocityUtil.generate(service_vm, service, context);
				System.out.println(service);
			}
			// 生成serviceMock
			File serviceMockFile = new File(serviceMock);
			if (!serviceMockFile.exists()) {
				VelocityContext context = new VelocityContext();
				context.put("package_name", packageName);
				context.put("model", model);
				context.put("ctime", ctime);
				VelocityUtil.generate(serviceMock_vm, serviceMock, context);
				System.out.println(serviceMock);
			}
			// 生成serviceImpl
			File serviceImplFile = new File(serviceImpl);
			if (!serviceImplFile.exists()) {
				VelocityContext context = new VelocityContext();
				context.put("package_name", packageName);
				context.put("model", model);
				context.put("model_first_low", StringUtil.toLowerCaseFirstOne(model));
				context.put("model_low", model.toLowerCase());
				context.put("ctime", ctime);
				VelocityUtil.generate(serviceImpl_vm, serviceImpl, context);
				System.out.println(serviceImpl);
			}

			//表字段包含id生成controller
			if(tables.get(i).get("has_autoKey").equals("true")){
				File controllerFile = new File(controller);
				if (!controllerFile.exists()) {
					VelocityContext context = new VelocityContext();
					context.put("package_name", packageName);
					context.put("model", model);
					context.put("model_first_low", StringUtil.toLowerCaseFirstOne(model));
					context.put("model_low", model.toLowerCase());
                    context.put("id", autoKey);
                    context.put("ctime", ctime);
					VelocityUtil.generate(controller_vm, controller, context);
					System.out.println(controller);
				}
			}
		}
		System.out.println("========== 结束生成Service ==========");
	}

	// 递归删除非空文件夹
	public static void deleteDir(File dir) {
		if (dir.isDirectory()) {
			File[] files = dir.listFiles();
			for (int i = 0; i < files.length; i++) {
				deleteDir(files[i]);
			}
		}
		dir.delete();
	}

}
